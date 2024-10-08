package config

import (
	"log"

	"github.com/bits-and-blooms/bloom/v3"
)

var EmailBloomFilter *bloom.BloomFilter

// Initialize Bloom Filter
// size: number of expected items to be added to the filter
// falsePositiveRate: acceptable false positive rate (between 0 and 1)
func InitBloomFilter(size uint, falsePositiveRate float64) {
	EmailBloomFilter = bloom.NewWithEstimates(size, falsePositiveRate)
	log.Println("Bloom filter initialized with size:", size, "and false positive rate:", falsePositiveRate)
}

// Add email to Bloom Filter
func AddEmailToBloom(email string) {
	EmailBloomFilter.AddString(email)
}

// Check if email might exist in Bloom Filter
func CheckEmailInBloom(email string) bool {
	return EmailBloomFilter.TestString(email)
}

// Custom hash function for Bloom Filter (optional)

// func fnvHash(data []byte) uint32 {
// 	hasher := fnv.New32()
// 	hasher.Write(data)
// 	return hasher.Sum32()
// }
