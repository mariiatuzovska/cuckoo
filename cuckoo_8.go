package cuckoo

import (
	"crypto/sha256"
	"errors"
	"math/rand"
)

// filter8 represents the structure of the filter with a byte fingerprint.
type filter8 struct {
	buckets     [][]byte
	size        int
	fingerprint func(key []byte) byte
}

// newCuckoo8Filter initializes a new Cuckoo Filter with a byte fingerprint.
func newCuckoo8Filter(numBuckets int) Filter {
	buckets := make([][]byte, numBuckets)
	for i := range buckets {
		buckets[i] = make([]byte, BucketSize*fingerprintSizeByte)
	}
	return &filter8{
		buckets:     buckets,
		fingerprint: func(key []byte) byte { return sha256.Sum256(key)[0] },
	}
}

// Insert adds an item to the Cuckoo Filter.
func (cf *filter8) Insert(key []byte) error {
	fp := cf.fingerprint(key)
	index1 := hash(key) % uint64(len(cf.buckets))

	// try to insert into either of the two buckets
	var added bool
	cf.buckets[index1], added = addFingerprint(cf.buckets[index1], fp)
	if !added {
		index2 := alternateIndex(index1, fp, len(cf.buckets))
		cf.buckets[index2], added = addFingerprint(cf.buckets[index2], fp)
	}
	if added {
		cf.size++
		return nil
	}

	// evict and relocate existing fingerprints if both buckets are full
	curIndex := index1
	for i := 0; i < MaxKicks; i++ {
		// randomly pick a fingerprint from the bucket to evict
		pos := rand.Intn(BucketSize)
		evictedFp := cf.buckets[curIndex][pos]
		cf.buckets[curIndex][pos] = fp

		// calculate the alternate index for the evicted fingerprint
		fp = evictedFp
		curIndex = alternateIndex(curIndex, fp, len(cf.buckets))
		cf.buckets[curIndex], added = addFingerprint(cf.buckets[curIndex], fp)
		if added {
			cf.size++
			return nil
		}
	}
	return errors.New("failed to insert after maximum number of kicks")
}

// Lookup checks if an item is in the Cuckoo Filter.
func (cf *filter8) Lookup(key []byte) bool {
	fp := cf.fingerprint(key)
	index1 := hash(key) % uint64(len(cf.buckets))
	index2 := alternateIndex(index1, fp, len(cf.buckets))

	// check if the fingerprint is in either bucket
	return containsFingerprint(cf.buckets[index1], fp) || containsFingerprint(cf.buckets[index2], fp)
}

// Delete removes an item from the Cuckoo Filter.
func (cf *filter8) Delete(key []byte) bool {
	fp := cf.fingerprint(key)
	index := hash(key) % uint64(len(cf.buckets))

	// try to remove the fingerprint from either bucket
	var removed bool
	cf.buckets[index], removed = removeFingerprint(cf.buckets[index], fp)
	if !removed {
		index = alternateIndex(index, fp, len(cf.buckets))
		cf.buckets[index], removed = removeFingerprint(cf.buckets[index], fp)
	}
	if removed {
		cf.size--
		return true
	}
	return false
}

// Size returns the number of elements in the Cuckoo Filter.
func (cf *filter8) Size() int {
	return cf.size
}
