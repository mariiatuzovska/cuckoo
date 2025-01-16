package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"time"

	"github.com/mariiatuzovska/cuckoo"
)

func main() {
	fmt.Println("---------- Cuckoo Filter with byte fingerprint ----------")
	cf, _ := cuckoo.NewCuckooFilter(cuckoo.FingerprintTypeByte, 100000)
	stats(cf, 64, 5000000, 100000)

	fmt.Println("---------- Cuckoo Filter with uint16 fingerprint ----------")
	cf, _ = cuckoo.NewCuckooFilter(cuckoo.FingerprintTypeUint16, 100000)
	stats(cf, 64, 5000000, 100000)

	fmt.Println("---------- Cuckoo Filter with uint32 fingerprint ----------")
	cf, _ = cuckoo.NewCuckooFilter(cuckoo.FingerprintTypeUint32, 100000)
	stats(cf, 64, 5000000, 100000)
}

func generateRandomStrings(count int, length int) [][]byte {
	strings := make([][]byte, count)
	for i := 0; i < count; i++ {
		b := make([]byte, length)
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		strings[i] = b
	}
	return strings
}

func stats(cf cuckoo.Filter, strSizeBytes, totalStrNum, checkSize int) {
	existingStrings := generateRandomStrings(totalStrNum, strSizeBytes)
	for _, s := range existingStrings {
		err := cf.Insert(s)
		if err != nil {
			panic(err)
		}
	}

	nonExistingStrings := generateRandomStrings(checkSize, strSizeBytes)

	// Measure memory usage after insertions
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Memory used after generating strings: %d KB\n", memStats.Alloc/1024)

	start := time.Now()

	// Check existing strings
	existingChecks := checkSize
	falseNegatives := 0
	for i := 0; i < existingChecks; i++ {
		if !cf.Lookup(existingStrings[i]) {
			falseNegatives++
		}
	}
	fmt.Printf("False negatives: %d (should be 0)\n", falseNegatives)

	// Check non-existing strings
	falsePositives := 0
	for _, s := range nonExistingStrings {
		if cf.Lookup([]byte(s)) {
			falsePositives++
		}
	}

	// Calculate and print error rate
	errorRate := float64(falsePositives) / float64(checkSize)
	fmt.Printf("False positives: %d\n", falsePositives)
	fmt.Printf("Error rate: %.6f\n", errorRate)

	// Print memory usage again after checks
	runtime.GC()
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Memory used after all operations: %d KB\n", memStats.Alloc/1024)

	// Measure and print processing time
	elapsed := time.Since(start)
	fmt.Printf("Processing time: %s\n", elapsed)
}
