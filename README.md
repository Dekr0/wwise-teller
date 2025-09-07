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

## Baseline Milestone

- Milestone 1: parse BKHD and skip all other chunks after reading chunk tag and 
chunk size.
    - Buffer Read: 20.92 ms
    - In Memory: 10.34 ms 
- Milestone 2: parse BKHD and enable parallel skip
