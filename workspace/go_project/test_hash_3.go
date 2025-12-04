package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

func getHashStr(tokenIds []uint32, priorHash string) string {
	hasher := sha256.New()
	
	if priorHash != "" {
		priorHashBytes, err := hex.DecodeString(priorHash)
		if err != nil {
			panic(fmt.Sprintf("Invalid prior hash: %v", err))
		}
		hasher.Write(priorHashBytes)
	}
	
	for _, t := range tokenIds {
		// Hash each token as 4-byte little-endian integer
		tokenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(tokenBytes, t)
		hasher.Write(tokenBytes)
	}
	
	return hex.EncodeToString(hasher.Sum(nil)) // Returns SHA256 hex string
}

func hashStrToInt64(hashStr string) int64 {
	// Take first 16 hex chars to get 64-bit value
	if len(hashStr) < 16 {
		panic("Hash string too short")
	}
	uint64Val, err := hex.DecodeString(hashStr[:16])
	if err != nil {
		panic(fmt.Sprintf("Invalid hex string: %v", err))
	}
	
	// Convert to uint64 (little-endian)
	val := binary.LittleEndian.Uint64(uint64Val)
	
	// Convert to signed int64 range [-2^63, 2^63-1]
	if val >= math.MaxInt64+1 { // 2^63
		return int64(val - math.MaxUint64)
	}
	return int64(val)
}

func main() {
	// SGLang's hash computation for first block
	tokens := []uint32{151644, 872, 198, 40, 1079, 730, 9991, 573, 23958, 11, 458, 20443, 11229, 17847, 7881, 553, 730, 4284, 812, 283, 13, 358, 1079, 6188, 311, 7789, 3847, 304, 35764, 4755, 11, 23163, 1467, 1741, 438, 7343, 11, 3946, 9293, 11, 14298, 11, 19502, 11, 19819, 32711, 11, 15473, 11, 323, 803, 11, 438, 1632, 438, 36710, 17979, 323, 5619, 3868, 13, 358, 15218, 3746}
	
	// Note: Go doesn't have a direct equivalent to Python's hash(tuple)
	// The Python line: print(hash(tuple(tokens)) + (1 << 64)) 
	// would need a different implementation in Go
	
	parentHash := ""
	hashVal := getHashStr(tokens, parentHash)
	
	fmt.Printf("Hash value: %s\n", hashVal)
	
	int64Val := hashStrToInt64(hashVal)
	fmt.Printf("Int64 value: %d\n", int64Val)
}
