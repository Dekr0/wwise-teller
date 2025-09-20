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
- Step 02:
    - Rethink about design:
        - Establish boundary between different data,
        - Completely separate data from behavior and logic, 
        - Structure data based on access pattern, efficient memory layout, and 
        boundary.

## Baseline Milestone

- Delta is measured against raw read using matching receiving buffer and reader 
buffer with size of 32k, 6.401 +/- 1%.
- Milestone 1: parse BKHD and skip all other chunks after reading chunk tag and 
chunk size.
    - Buffer size 32k: 6.741 ms +/- 3%
    - Delta: 0.34 ms +/- 3%
- Milestone 2: parse BKHD and store encoded chunk
    - Buffer size 32k: 9.376 ms +/- 3%
    - Delta: 2.975 ms +/- 0%
- Milestone 3: Milestone 2 + parse DIDX
    - Buffer size 32k: 11.62 ms +/- 0%
    - Delta: 5.219 ms +/- 0%
- Milestone 4: Milestone 3 + parse DATA
    - Buffer size 32k: 11.77 ms +/- 1%
    - Delta: 5.369 ms +/- 0%
- Milestone 5: Concurrent parsing (parse DIDX)
    - Buffer size 32k, 4 chunk decoder: 10.34 +/- 1%
        - If buffer size is greater than chunk size, use chunk size * 2
    - Delta: 3.939 ms +/- 0%
- Milestone 5: Milestone 5 + parse DATA
    - Buffer size 32k, 4 chunk decoder: 10.74 +/- 0%
        - If buffer size is greater than chunk size, use chunk size * 2
    - Delta: 4.339 ms +/- 1%
