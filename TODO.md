# Priority List

- The ability to modify DATA and DIDX section
    - Automation as first class citizen 
- The ability to pack multiple sound banks for Helldivers 2 integration

# Uncategorized Feature Idea

- Attempt to add new Sound object in the hierarchy
    - If this work, this mean hierarchy modification is completely possible
- Information copying
- Show all properties (read only)
- Modify the destination file name
- Ability to delete / rename file in file explorer and file dialog
- Reverse engineering on `wwise_properties` and `wwise_metadata` that appears in 
different archive file.
- Reverse engineering on the wwise configuration and asset path that appear in 
the `setting.ini` file.
- Experiment the possibility of packing sound bank without its Wwise dependency

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

- Command palette should contain "save sound banks" option
    - It might be better put them into command palette that is dedicated for 
    bank explorer.
- For file explorer and file dialog, split a file path into individual segment.
When clicking on a segment, it will jump to a specific file location.
- Overhaul configuration setting

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

## Encoding

- There seems to be a lot of heap allocations during encoding phase. HIRC encoding 
is the most.
- Why there's a lot of heap allocations? Each encoding thread allocate memory 
to fill up the encoded data, and return that encoded data to the main thread.
This is done for laziness because I don't need to calculate and determine 
the section of memory a thread can write to.

### Rework Proposal

- The lifetime of encoding data is known. It ends once the encoded data write 
into the disk.
- Size of every chunk is known.
- Caller ask the size of a bank, and uses this size to allocate a fix chunk of 
memory upfront.
- The bank encoding function use this memory.

#### Difficulty

- Each thread needs to know the start position and end position in this memory 
block so that it doesn't write into the area of other threads.

#### Solution

- Before scheduling a thread for encoding something, main thread caluclate the 
start position and end position for this thread, and create a section writer that 
keep track of cursor using these two positions so that this thread doesn't write 
out of bound, or doesn't fill up the entire section as requested.
- This section writer will immedately panic if either of this situation happen 
since encoding must not fail.

### Potential Benefit

- Memory recycle become extremely trivial since there's only a block. Once the 
encoded data is write to the disk., it can be recycled.
- Easy to detect any encoding logic bug, a single byte off will catch by the 
panic / assertion of the section writer.

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
