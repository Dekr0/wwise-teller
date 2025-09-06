# Wwise Teller

- This branch is a completely rework and redesign of Wwise Teller. This branch 
will attempt address all issues and mistakes in the `main` branch.
- The choice of programming language for rework of Wwise Teller will still be 
Go. There are few reasons:
    - I want to continue to improve my understanding and usage with Go. This 
    will inclueds:
        - Potential pitfalls (e.g., anti-pattern, unexpected behavior, etc.)
        - Performance technique for working with a GC language
    - I originally want to use either `Zig`, `Odin`, or `C` for the rework. 
    Wwise is a fairly complex sound engine, and there are still decent amount 
    of domains I have yet scoped out. Another thing is that I have yet spent 
    enough time with these languages I mentioned. Thus, I will stick with Go at 
    the moment since I can make it reasonly performant while have some amount of 
    headrooms for prototypes or hitting a design flaw.
- The `main` branch will remain active but it will primarly serve as a branch 
for prototype and exploring different behavior in Wwise.

## Rework Roadmap

- Step 01: 
    - IO function and utility function rework
    - Establish performance and implementation baseline
    - DON'T LOST THE BASELINE. Try to keep it as close as possible

### Baseline

#### IO Speed of Development Hard Disk

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

#### IO Speed of Average Consumer Hard Disk

#### Fastest IO speed in Go

- Cold load
- Hot load (OS optimization primitive such as file system cache is active)
