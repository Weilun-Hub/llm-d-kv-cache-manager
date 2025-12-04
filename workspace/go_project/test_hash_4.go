package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Python风格的tuple哈希实现
func hashTuple(tokens []uint32) int64 {
	// Python的tuple哈希算法：对每个元素递归哈希并组合
	// 简化版：使用FNV-like算法模拟Python行为
	
	const (
		mult = 1000003
		xor  = 0x345678
	)
	
	hashVal := int64(xor)
	
	for _, token := range tokens {
		// 对每个token进行哈希（Python对int的哈希就是值本身）
		tokenHash := int64(token)
		
		// Python的tuple哈希组合公式
		hashVal = (hashVal ^ tokenHash) * mult
		
		// 防止溢出，模拟Python的整数回绕
		hashVal &= math.MaxInt64
	}
	
	// 添加tuple长度的影响
	hashVal += int64(len(tokens))
	hashVal &= math.MaxInt64
	
	// Python的哈希需要是有符号的
	if hashVal > math.MaxInt32 {
		hashVal -= math.MaxInt32 * 2
	}
	
	return hashVal
}

// 更精确的实现：使用与Python相同的算法
func pythonHashTuple(tokens []uint32) int64 {
	// 这个实现更接近Python的实际算法
	var h int64 = 0x345678
	
	for _, x := range tokens {
		// Python对int的哈希：对于小整数就是值本身
		// 对于大整数有更复杂的算法，但这里token ID通常不大
		h = (h ^ int64(x)) * 1000003
		h &= 0xFFFFFFFFFFFFFFFF // 64位掩码
	}
	
	h ^= int64(len(tokens))
	
	// 转换为有符号64位整数
	if h > math.MaxInt64 {
		h = -((^h) + 1) // 二进制补码
	}
	
	return h
}

func main() {
	tokens := []uint32{151644, 872, 198, 40, 1079, 730, 9991, 573, 23958, 11, 458, 20443, 11229, 17847, 7881, 553, 730, 4284, 812, 283, 13, 358, 1079, 6188, 311, 7789, 3847, 304, 35764, 4755, 11, 23163, 1467, 1741, 438, 7343, 11, 3946, 9293, 11, 14298, 11, 19502, 11, 19819, 32711, 11, 15473, 11, 323, 803, 11, 438, 1632, 438, 36710, 17979, 323, 5619, 3868, 13, 358, 15218, 3746}
	
	// 计算Python风格的tuple哈希
	pythonHash := pythonHashTuple(tokens)
	fmt.Printf("Python hash(tuple(tokens)): %d\n", pythonHash)
	
	// 加上 (1 << 64)
    // result := pythonHash + (1 << 64)
    // fmt.Printf("hash(tuple(tokens)) + (1 << 64): %d\n", result)
}

