package main

import (
    "encoding/binary"
    "fmt"
    "hash/fnv"
)

// computePythonHash computes a hash similar to Python's hash(tuple(tokens)) and returns uint64
func computePythonHash(tokens []uint32) uint64 {
    // Use FNV-1a hash as a simple approximation
    hasher := fnv.New64a()
    for _, token := range tokens {
        tokenBytes := make([]byte, 4)
        binary.LittleEndian.PutUint32(tokenBytes, token)
        hasher.Write(tokenBytes)
    }
    // Return as uint64 directly (Go's uint64 conversion handles two's complement)
    return hasher.Sum64()
}

func main() {
    // Test tokens from SGLang event
    tokens := []uint32{151644, 872, 198, 40, 1079, 730, 9991, 573, 23958, 11, 458, 20443, 11229, 17847, 7881, 553, 730, 4284, 812, 283, 13, 358, 1079, 6188, 311, 7789, 3847, 304, 35764, 4755, 11, 23163, 1467, 1741, 438, 7343, 11, 3946, 9293, 11, 14298, 11, 19502, 11, 19819, 32711, 11, 15473, 11, 323, 803, 11, 438, 1632, 438, 36710, 17979, 323, 5619, 3868, 13, 358, 15218, 3746}

    hashUint64 := computePythonHash(tokens)
    
    fmt.Printf("Tokens: %v\n", tokens)
    fmt.Printf("Go computePythonHash result: %d\n", hashUint64)
    fmt.Printf("Go computePythonHash (hex): 0x%x\n", hashUint64)
    
    // Convert to int64 to see if it would be negative in Python
    hashInt64 := int64(hashUint64)
    fmt.Printf("As int64: %d\n", hashInt64)
    
    // Expected from Python: hash(tuple(tokens)) + (1 << 64) if negative
    // Or just hash(tuple(tokens)) converted to uint64
    fmt.Println("\nCompare with Python:")
    fmt.Println("  Python: hash(tuple(tokens))")
    fmt.Println("  If negative, convert: hash(tuple(tokens)) + (1 << 64)")
    fmt.Println("  Or use: ctypes.c_uint64(hash(tuple(tokens))).value")
}

