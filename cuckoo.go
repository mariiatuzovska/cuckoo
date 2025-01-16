package cuckoo

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// Filter represents the interface for a Cuckoo Filter.
type (
	Filter interface {
		Insert(key []byte) error
		Lookup(key []byte) bool
		Delete(key []byte) bool
		Size() int
	}

	FingerprintType int
)

const (
	// Base types for fingerprints
	FingerprintTypeByte FingerprintType = iota // Type of fingerprint in bytes.
	FingerprintTypeUint16
	FingerprintTypeUint32

	BucketSize = 256  // Number of fingerprints per bucket
	MaxKicks   = 1000 // Max re-insertions before failing

	// Size of fingerprint in bytes.
	fingerprintSizeByte   = 1
	fingerprintSizeUint16 = 2
	fingerprintSizeUint32 = 4
)

func NewCuckooFilter(t FingerprintType, numBuckets int) (Filter, error) {
	switch t {
	case FingerprintTypeByte:
		return newCuckoo8Filter(numBuckets), nil
	case FingerprintTypeUint16:
		return newCuckoo16Filter(numBuckets), nil
	case FingerprintTypeUint32:
		return newCuckoo32Filter(numBuckets), nil
	}
	return nil, errors.New("unsupported fingerprint type")
}

// hash computes the hash of a given key.
func hash(data []byte) uint64 {
	h := sha256.Sum256(data)
	return binary.BigEndian.Uint64(h[:8])
}

// alternateIndex calculates the alternate index for a fingerprint.
func alternateIndex[T ~byte | ~uint16 | ~uint32 | ~uint64](index uint64, fp T, size int) uint64 {
	return (index ^ uint64(fp)) % uint64(size)
}

// removeFingerprint removes a fingerprint from a bucket.
func removeFingerprint[T ~byte | ~uint16 | ~uint32 | ~uint64](bucket []T, fp T) ([]T, bool) {
	for i := 0; i < BucketSize; i++ {
		if bucket[i] == fp {
			bucket[i] = 0
			return bucket, true
		}
	}
	return bucket, false
}

// addFingerprint adds a fingerprint to a bucket if there is space.
func addFingerprint[T ~byte | ~uint16 | ~uint32 | ~uint64](bucket []T, fp T) ([]T, bool) {
	for i := 0; i < BucketSize; i++ {
		if bucket[i] == 0 {
			bucket[i] = fp
			return bucket, true
		}
	}
	return bucket, false
}

// containsFingerprint checks if a fingerprint is in a bucket.
func containsFingerprint[T ~byte | ~uint16 | ~uint32 | ~uint64](bucket []T, fp T) bool {
	for i := 0; i < BucketSize; i++ {
		if bucket[i] == fp {
			return true
		}
	}
	return false
}
