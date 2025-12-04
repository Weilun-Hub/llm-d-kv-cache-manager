package main

import (
    "fmt"
)

// Python's hash constants (from sys.hash_info)
const (
    // Python 3 uses SipHash24, but for tuple hashing, it uses a simpler algorithm
    // Tuple hash: acc = (acc * multiplier + hash(elem)) % modulus
    // For integers, hash(i) = i (for small integers)
    // Initial value for tuple hash
    pythonTupleHashInit = 0x345678
    // Multiplier used in tuple hash combination
    // Python uses different multipliers, but for simplicity we'll use a common one
    pythonHashMultiplier = 1000003
    // Modulus: 2^61 - 1 (Python's default)
    pythonHashModulus = 2305843009213693951
)

// hashInt computes Python's hash for an integer
// For small integers, Python's hash(i) = i
func hashInt(val uint32) int64 {
    return int64(val)
}

// computePythonHash computes Python's hash(tuple(tokens)) and returns uint64
// Python's tuple hash algorithm:
//   acc = initial_value
//   for each element:
//     acc = (acc * multiplier + hash(element)) % modulus
//   return acc (which can be negative, so we convert to uint64)
func computePythonHash(tokens []uint32) uint64 {
    acc := int64(pythonTupleHashInit)
    
    for _, token := range tokens {
        elemHash := hashInt(token)
        // acc = (acc * multiplier + elemHash) % modulus
        // Use int64 arithmetic, handling overflow
        acc = ((acc * pythonHashMultiplier) + elemHash) % pythonHashModulus
    }
    
    // Python's hash can be negative, convert to uint64
    // Go's uint64() conversion automatically handles two's complement:
    // negative int64 values are converted to uint64 by reinterpreting the bits
    return uint64(acc)
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

