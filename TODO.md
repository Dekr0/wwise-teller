# Performance

## Problems

## File IO

- The hierarchy decoding process duplicates portion of bytes from the entire 
bytes slices of hierarchy section to avoid race condition. This should be 
done through partition using index marking so that each decoding go routine does 
not interfere each other.

## Go Routine

- There is no worker pool for hierarchy decoding go routines and hierarchy 
encoding go routines.

# Bank Table Selection

- I want to maintain what I already selected despite filter is applied
- Bug:
    - selection stage in object_editor
