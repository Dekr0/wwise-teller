# Performance

## Problems

## File IO

- The hierarchy decoding process duplicates portion of bytes from the entire 
bytes slices of hierarchy section to avoid race condition. 

## Go Routine

- There is no worker pool for hierarchy decoding go routines and hierarchy 
encoding go routines.

# UI & UX 

- Maintain what is being selected despite filter is applied
