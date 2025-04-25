# Design Dump

## Performance

### Pitfall

#### Page Fault

- Reuse memory buffer that is used for reading file to prevent page fault?

#### Memory Duplication

- At decoding phase, there are memory duplication. This is intended at this stage 
because it prevent multiple go routine share the exact same reader, specifically, 
modify the cursor position of the share reader.
    - Solution, create a new type of reader where it operates the same buffer but 
    it maintains its own cursor for each instance of the reader.

### Bottleneck

#### Parallelism

- The current implementation of parallel parsing will bottle neck on the following 
factors:
    1. How fast main routine can collect the result from parser routine and handle 
    bookkeeping
    2. Channel contention. If I understand it right, channel can be viewed as a 
    C thread safe queue (without buffered channel) or circular buffer with 
    conditioning (buffered channel). So, there will be contention when parser 
    routine try to write into the channel.
    3. Parser complexity
