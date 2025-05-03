# Performance

## Problems

## File IO

- The hierarchy decoding process duplicates portion of bytes from the entire 
bytes slices of hierarchy section to avoid race condition. 

## Go Routine

- There is no worker pool for hierarchy decoding go routines and hierarchy 
encoding go routines.

# UI & UX 

- Command palette should contain "save sound banks" option
    - It might be better put them into command palette that is dedicated for 
    bank explorer.
- Maintain what is being selected despite filter is applied
- Input Text for config
- Input Bugs for range properties

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
