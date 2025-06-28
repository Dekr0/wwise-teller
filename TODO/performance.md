# Struct Packing and Space Efficiency

- Hierarchy `struct` should use value and value of `struct` instead of 
pointer to `struct` and array of pointers to `struct`

# File IO

- Buffer read
    - Use multiple file descriptor when multi-threading
    - Limit multi-threading when read involve disk access
- The hierarchy decoding process duplicates portion of bytes from the entire 
bytes slices of hierarchy section to avoid race condition. 

# Encoding

- There seems to be a lot of heap allocations during encoding phase. HIRC encoding 
is the most.
- Why there's a lot of heap allocations? Each encoding thread allocate memory 
to fill up the encoded data, and return that encoded data to the main thread.
This is done for laziness because I don't need to calculate and determine 
the section of memory a thread can write to.

## Rework Proposal

- The lifetime of encoding data is known. It ends once the encoded data write 
into the disk.
- Size of every chunk is known.
- Caller ask the size of a bank, and uses this size to allocate a fix chunk of 
memory upfront.
- The bank encoding function use this memory.

### Difficulty

- Each thread needs to know the start position and end position in this memory 
block so that it doesn't write into the area of other threads.

### Solution

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

# Go Routine

- There is no worker pool for hierarchy decoding go routines and hierarchy 
encoding go routines.
