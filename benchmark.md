# Baseline

## IO speed Of development storage 

- SN750 500 GB
- Datasheet
    - PCIe Gen4
    - Sequential Read: 3.6GB/s
    - Sequential Write: 3.6GB/s
    - Random Read: 360k IOPS
    - Random Write: 480k IOPS
- CrystalDiskMark
    - SEQ1M (Q8T1 | 4 Test Counts | 2Gib Test Size)
        - 3.4925 GB/s Read 
        - 2.60260 GB/s Write

## Raw read speed in Go using developement hard disk

- Benchmark data from `go test -bench` has a very high chance of which it 
reflects the situation of hot loaded.
- Time is measured in milisecond.

### Matching input rate and output rate

- For buffer read, buffer size is the same as the destination buffer size.

- Total size is 36282368 bytes.
- Read all at once -> 10.70 ms +/- 10%
```
----------------------------------------
| size | unbuffer     | buffer         |
|------|--------------|----------------|
| 4k   | 20.72 +/- 2% | 20.14 +/- 2%   |
| 8k   | 12.51 +/- 4% | 12.76 +/- 2%   |
| 16k  | 8.664 +/- 7% | 8.389 +/- 3%   |
| 32k  | 6.256 +/- 0% | 6.286 +/- 0%   |
| 32k  | 5.523 +/- 3% | 5.619 +/- 12%  |
| 32k  | 5.084 +/- 1% | 5.124 +/- 2%   |
----------------------------------------
```

- Total size is 1758597120 bytes
- Read all at once -> 539.8 +/- 1%
```
-------------------------
| size  | buffer        |
|-------|---------------|
| 4k    | 992.0 +/- 2%  |
| 8k    | 642.1 +/- 3%  |
| 16k   | 430.4 +/- 1%  |
| 32k   | 314.4 +/- 0%  |
| 64k   | 272.7 +/- 1%  |
| 128k  | 260.1 +/- 6%  |
-------------------------
```

### Mismatch input rate and output rate

- Buffer size is the different than the destination buffer size.

- Total size is 36282368 bytes.
- Read all at once -> 11.37
```
------------------------------------------------------
| recv buffer size | read buffer size | time         |
|------------------|------------------|--------------|
| 1k               | 32k              | 347.3 +/- 1% |
| 2k               | 32k              | 342.7 +/- 2% |
| 4k               | 32k              | 345.4 +/- 1% |
| 8k               | 32k              | 347.2 +/- 1% |
| 16k              | 32k              | 353.7 +/- 2% |
| 32k              | 32k              | 319.3 +/- 1% |
------------------------------------------------------
```

- Total size is 1758597120 bytes
- Read all at once -> 517.9
```
------------------------------------------------------
| recv buffer size | read buffer size | time         |
|------------------|------------------|--------------|
| 1k               | 32k              | 6.989 +/- 3% |
| 2k               | 32k              | 6.896 +/- 3% |
| 4k               | 32k              | 6.911 +/- 3% |
| 8k               | 32k              | 7.090 +/- 3% |
| 16k              | 32k              | 7.044 +/- 8% |
| 32k              | 32k              | 6.401 +/- 1% |
------------------------------------------------------
```

## Raw write speed in Go using developement hard disk

### Matching input rate and output rate

- Total size is 36282368 bytes.
- Write all at once -> 20.83 ms +/- 4%
```
------------------------------
| Write size | Time          |
|------------|---------------|
| 4k         | 40.87 +/- 13% |
| 8k         | 25.06 +/- 3%  |
| 16k        | 18.02 +/- 2%  |
| 32k        | 15.28 +/- 1%  |
| 64k        | 13.80 +/- 2%  |
| 128k       | 12.74 +/- 1%  |
------------------------------
```
