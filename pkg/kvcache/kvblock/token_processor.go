/*
Copyright 2025 The llm-d Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kvblock

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
    "encoding/hex"
    "fmt"
    "strconv"
    "hash/fnv"

	"github.com/fxamacker/cbor/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/llm-d/llm-d-kv-cache-manager/pkg/utils"
	// "github.com/llm-d/llm-d-kv-cache-manager/pkg/utils/logging"

)

// defaultBlockSize is the default number of tokens per block.
// 16 is the default value used by vLLM.
const defaultBlockSize = 16

// TokenProcessorConfig holds the configuration for the token processor.
type TokenProcessorConfig struct {
	BlockSize int `json:"blockSize"`
	// HashSeed is used to prefix initial hash chunks, similarly to vLLM's NONE_HASH.
	// This should be aligned with vLLM's `PYTHONHASHSEED` environment variable.
	// The system's deployer is responsible for aligning the vLLM deployments
	// with the same seed value.
	HashSeed string `json:"hashSeed"`
	initHash []byte // cache once
}

// DefaultTokenProcessorConfig returns the default configuration for the token processor.
func DefaultTokenProcessorConfig() *TokenProcessorConfig {
	return &TokenProcessorConfig{
		BlockSize: defaultBlockSize,
		HashSeed:  "",
	}
}

// TokenProcessor defines the interface for converting tokens to
// KVBlockKeys.
type TokenProcessor interface {
	// TokensToKVBlockKeys converts tokens into kv_block.Keys.
	TokensToKVBlockKeys(tokens []uint32, modelName string) []Key
	TokensToKVBlockKeysSGLang(tokens []uint32, modelName string, parentHashHex string) []Key
}

// ChunkedTokenDatabase is a concrete implementation of TokenDatabase.
// It mimics the ChunkedTokenDatabase in the Python code.
type ChunkedTokenDatabase struct {
	TokenProcessorConfig
}

var _ TokenProcessor = &ChunkedTokenDatabase{}

// NewChunkedTokenDatabase creates a new instance with the given config and metadata.
func NewChunkedTokenDatabase(config *TokenProcessorConfig) TokenProcessor {
	if config == nil {
		config = DefaultTokenProcessorConfig()
	} // TODO: validate?

	return &ChunkedTokenDatabase{
		TokenProcessorConfig: *config,
	}
}

// getInitHash returns the root parent hash as a full byte slice.
func (db *ChunkedTokenDatabase) getInitHash() []byte {
	if db.initHash != nil {
		return db.initHash
	}

	encMode, err := cbor.CanonicalEncOptions().EncMode() // deterministic
	if err != nil {
		log.FromContext(context.Background()).Error(err, "failed to create CBOR encoder")
		return nil
	}

	b, err := encMode.Marshal(db.HashSeed)
	if err != nil {
		log.FromContext(context.Background()).Error(err, "failed to marshal payload to CBOR")
		return nil
	}

	sum := sha256.Sum256(b)
	db.initHash = sum[:] // Return the full 32-byte hash
	return db.initHash
}

func HashUint32Tokens(tokens []uint32) uint64 {
    h := fnv.New64a()
    var buf [4]byte

    for _, t := range tokens {
        tokenBytes := make([]byte, 4)
        binary.LittleEndian.PutUint32(tokenBytes, t)
        h.Write(tokenBytes)
    }

    hashVal := int64(h.Sum64())
    if hashVal >

    return h.Sum64()
}

// hash computes the full 32-byte SHA256 hash of the given parent, tokens,
// and extra keys, mimicking the vLLM implementation.
func (db *ChunkedTokenDatabase) hash(parent []byte, tokens []uint32, extra interface{}) []byte {
	payload := []interface{}{parent, tokens, extra}

	encMode, err := cbor.CanonicalEncOptions().EncMode() // deterministic
	if err != nil {
		log.FromContext(context.Background()).Error(err, "failed to create CBOR encoder")
		return nil
	}

	b, err := encMode.Marshal(payload)
	if err != nil {
		log.FromContext(context.Background()).Error(err, "failed to marshal payload to CBOR")
		return nil
	}

	sum := sha256.Sum256(b)
	return sum[:] // Return the full 32-byte hash
}

func (db *ChunkedTokenDatabase) hashSGLang(parentHashHex string, tokens []uint32) string {
    hasher := sha256.New()

    if parentHashHex != "" {
        parentBytes, err := hex.DecodeString(parentHashHex)
        if err == nil {
            hasher.Write(parentBytes)
        }
    }
    
    tokenBytes := make([]byte, 4)
    for _, token := range tokens {
        binary.LittleEndian.PutUint32(tokenBytes, token)
        hasher.Write(tokenBytes)
    }

    return fmt.Sprintf("%x", hasher.Sum(nil))
}

func hashStrToInt64(hashStr string) int64 {
    if len(hashStr) < 16 {
        return 0
    }

    uint64Val, err := strconv.ParseUint(hashStr[:16], 16, 64)
    if err != nil {
        return 0
    }

    // uint64Val := uint64(0)
    // fmtScanf(hashStr[:16], "%x", &uint64Val)
    // const maxInt64Plus1 = uint64(1) << 63
    // if uint64Val >= maxInt64Plus1 {
    //    return int64(uint64Val - maxInt64Plus1) - int64(maxInt64Plus1)
    //}

    return int64(uint64Val)
}

func (db *ChunkedTokenDatabase) TokensToKVBlockKeysSGLang(tokens []uint32, modelName string, parentHashHex string) []Key {
    chunks := db.chunkTokens(tokens)
    if len(chunks) == 0 {
        return nil
    }

    keys := make([]Key, 0, len(chunks))
    parentHash := parentHashHex

    logger := log.FromContext(context.Background())
    logger.Info("TokensToKVBlockKeysSGLang starting",
        "parentHashHex", parentHashHex,
        "chunkCount", len(chunks),
        "firstChunkTokens", chunks[0])

    for i, chunk := range chunks {
        hashHex := db.hashSGLang(parentHash, chunk)
        hashInt64 := hashStrToInt64(hashHex)

        //var hashUint64 uint64
        //if hashInt64 < 0 {
        //    const maxInt64Plus1 = uint64(1) << 63
        //    hashUint64 = uint64(hashInt64) + maxInt64Plus1 + maxInt64Plus1
        //} else {
        //    hashUint64 = uint64(hashInt64)
        //}
        //hashUint64 := uint64(hashInt64)

        hashUint64 := HashUint32Tokens(chunk)

        logger.Info("SGLang hash computation",
            "chunkIndex", i,
            "parentHashHex", parentHash,
            "chunkTokens", chunk,
            "fullHashHex", hashHex,
            "hashInt64", hashInt64,
            "hashUint64", hashUint64)

        keys = append(keys, Key{
            ModelName: modelName,
            ChunkHash: hashUint64,
        })

        parentHash = hashHex
    }

    return keys

}
// prefixHashes returns a slice of full 32-byte hashes.
func (db *ChunkedTokenDatabase) prefixHashes(parentHash []byte, tokenChunks [][]uint32) [][]byte {
	prefix := parentHash
	hashes := make([][]byte, len(tokenChunks))
	for i, chunk := range tokenChunks {
		prefix = db.hash(prefix, chunk, nil)
		hashes[i] = prefix
	}
	return hashes
}

// chunkTokens splits the input slice of tokens into chunks of size chunkSize.
func (db *ChunkedTokenDatabase) chunkTokens(tokens []uint32) [][]uint32 {
	var chunks [][]uint32
	for i := 0; i < len(tokens); i += db.BlockSize {
		end := i + db.BlockSize
		if end > len(tokens) {
			break // no partial blocks
		}

		chunks = append(chunks, tokens[i:end])
	}

    if len(chunks) > 0 {
        log.FromContext(context.Background()).Info("chunked tokens",
            "blockSize", db.BlockSize,
            "totalTokens", len(tokens),
            "chunkCount", len(chunks),
            "firstChunkTokens", chunks[0],
            "hashSeed", db.HashSeed)
    }

	return chunks
}

// TokensToKVBlockKeys converts tokens into kv_block.Keys.
func (db *ChunkedTokenDatabase) TokensToKVBlockKeys(tokens []uint32, modelName string) []Key {
	parentBytes := db.getInitHash()
	if parentBytes == nil {
		return nil
	}

    initHashHex := ""
    if len(parentBytes) >= 8 {
        initHashHex = fmt.Sprintf("%x", parentBytes[:8])
    }
    log.FromContext(context.Background()).Info("generating block keys",
        "hashSeed", db.HashSeed,
        "initHashPrefix", initHashHex,
        "blockSize", db.BlockSize)

	chunks := db.chunkTokens(tokens)
	ph := db.prefixHashes(parentBytes, chunks)

    computedHashesInfo := make([]string, 0, len(ph))
    computedUint64Hashes := make([]uint64, 0, len(ph))
    for i, hashBytes := range ph {
        hashHex := fmt.Sprintf("%x", hashBytes)
        hashVal := binary.BigEndian.Uint64(hashBytes[24:])
        computedHashesInfo = append(computedHashesInfo, fmt.Sprintf("chunk[%d]: full=%s uint64=%d", i, hashHex, hashVal))
        computedUint64Hashes = append(computedUint64Hashes, hashVal)
    }
    log.FromContext(context.Background()).Info("KV manager computed block hashes",
        "computedHashes", computedHashesInfo,
        "computedUint64Hashes", computedUint64Hashes,
        "chunkCount", len(chunks))

	// Convert the final byte hashes to uint64 for the Key struct
	return utils.SliceMap(ph, func(hashBytes []byte) Key {
		// Truncate to 64 bits at the very end by taking the last 8 bytes
		hashVal := binary.BigEndian.Uint64(hashBytes[24:])
		return Key{
			ModelName: modelName,
			ChunkHash: hashVal,
		}
	})
}
