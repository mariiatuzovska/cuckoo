# Cuckoo Filter Implementation and Performance

This project demonstrates the implementation of a Cuckoo Filter with varying fingerprint sizes (`byte`, `uint16`, and `uint32`) to evaluate memory usage, error rates, and processing time for inserting and querying elements.

## Overview of Cuckoo Filter
A **Cuckoo Filter** is a probabilistic data structure used for set membership testing with efficient support for adding and removing elements. It is space-efficient but allows for a tunable false positive rate.

## Results Summary
The following table summarizes the performance metrics for different fingerprint sizes:

| Fingerprint Type | Memory Used After Generating Strings (KB) | False Negatives | False Positives | Error Rate | Memory Used After All Operations (KB) | Processing Time |
|------------------|-------------------------------------------|-----------------|-----------------|------------|---------------------------------------|-----------------|
| **Byte**         | 465,810                                   | 0               | 32,259          | 0.322590   | 182                                   | 85.110ms        |
| **Uint16**       | 540,826                                   | 0               | 150             | 0.001500   | 198                                   | 93.631ms        |
| **Uint32**       | 840,826                                   | 0               | 0               | 0.000000   | 198                                   | 107.181ms       |

## Parameters Used
The performance metrics above were generated using the following parameters:

- **Number of Buckets**: 100,000
- **String Size**: 64 bytes
- **Total Strings Inserted**: 5,000,000
- **Number of Checks (Existing and Non-Existing)**: 100,000

The `main` function used for these experiments can be found in [./cmd/main.go](./cmd/main.go):

## Key Observations
- **Memory Usage**:
  - Larger fingerprint sizes result in increased memory usage.
  - The `uint32` filter consumes significantly more memory compared to the `byte` filter.

- **False Positives**:
  - Smaller fingerprints have higher false positive rates.
  - The `byte` filter exhibits a noticeable error rate of 32.2%, whereas the `uint32` filter achieves 0 false positives.

- **Processing Time**:
  - Processing time increases with fingerprint size due to additional computations and memory overhead.

## How to Run

```bash
go run ./cmd/main.go
```

## Code Highlights
The project includes:
- **Implementation of Cuckoo Filters** with varying fingerprint sizes.
- Functions to measure performance metrics such as memory usage, false positive rate, and processing time.
- A unified interface for filter operations (`Insert`, `Lookup`, `Delete`, and `Size`).

## Future Improvements
- Support for dynamic resizing of filters.
- Integration with other hash functions for performance comparison.
- Enhanced benchmarking with larger datasets and more query types.
