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

## Fastest IO speed in Go using developement hard disk

- Benchmark data from `go test -bench` has a very high chance of which it 
reflects the situation of hot loaded.
- Time is measured in milisecond.

### Matching input rate and output rate

- For buffer read, buffer size is the same as the destination buffer size.

#### `weapon_superearth` sound bank

- Total size is 36282368 bytes.
- Read all at once -> 11.37

```
----------------------------
| size | unbuffer | buffer |
|------|----------|--------|
| 4k   | 21.45    | 21.34  |
| 8k   | 13.05    | 13.05  |
| 16k  | 9.047    | 9.057  |
| 32k  | 6.868    | 6.850  |
| 64k  | 5.910    | 5.944  |
| 128k | 5.741    | 5.769  |
----------------------------
```

#### Largest gaem archive in Helldivers 2

- Total size is 1758597120 bytes
- Read all at once -> 517.9

```
-----------------
| size | buffer |
|------|--------|
| 4k   | 1004   |
| 8k   | 607.2  |
| 16k  | 404.5  |
| 32k  | 306.4  |
| 64k  | 255.6  |
| 128k | 236.6  |
|---------------|

```

### Mismatch input rate and output rate

- Buffer size is the different than the destination buffer size.

#### `weapon_superearth` sound bank

- Total size is 36282368 bytes.
- Read all at once -> 11.37

```
-----------------------------------
| size | recv buffer | read buffer |
|------|-------------|-------------|
| 4k   | 21.45       | 21.34       |
| 8k   | 13.05       | 13.05       |
| 16k  | 9.047       | 9.057       |
| 32k  | 6.868       | 6.850       |
| 64k  | 5.910       | 5.944       |
| 128k | 5.741       | 5.769       |
------------------------------------
```
