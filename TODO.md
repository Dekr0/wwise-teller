# Features

## Rewiring

### Change parent of a hierarchy & Add / Remove child from a container

- Update: this only works for node with no leafs!

#### Change parent of a hierarchy

- Flow:
    1. The parent remove the target hierarchy.
    2. If step 1 success, the new parent add the target hierarchy.
    3. If step 2 success, change direct parent ID of the target hierarchy.
    4. Arrange the tree index. Put the target hierarchy next to the new parent

#### Remove a child from a container

- Flow:
    1. The parent remove the target hierarchy.
    2. If step 1 success, change direct parent ID of the target hierarchy to 
    nothing.
    3. Arrange the tree index. Put the target hierarchy to free float.

## UI & UX 

- Display audio source ID in playlist setting.
- Command palette should contain "save sound banks" option
    - It might be better put them into command palette that is dedicated for 
    bank explorer.
- For file explorer and file dialog, split a file path into individual segment.
When clicking on a segment, it will jump to a specific file location.
- Overhaul configuration setting
- Object editor requires Bank explore to present so it can display all selected 
hierarchy objects.

# Performance

## Problems

## Data Compact and Pointer

- Hierarchy `struct` should use value and value of `struct` instead of 
pointer to `struct` and array of pointers to `struct`

## File IO

- Buffer read
    - Use multiple file descriptor when multi-threading
    - Limit multi-threading when read involve disk access
- The hierarchy decoding process duplicates portion of bytes from the entire 
bytes slices of hierarchy section to avoid race condition. 

## Go Routine

- There is no worker pool for hierarchy decoding go routines and hierarchy 
encoding go routines.

# Design

- Callback chaining?
- Reactive primitives since state changes might get out hand. 
    - Fine grained synchronous reactivity
        - Signal
    - Or, Event Bus / Stream asynchronous reactivity
        - producer emit data into buses / streams
        - emit data into buses / streams
        - buses / streams transform into new types of buses / streams by applying 
        functional like primitives such as `map`, `buffer`, `reduce`, etc.
        - consumer consume data emitted from buses / streams

