# Design, Rule, and Convention

## Panic vs. Error Return

- All hierarchy decoding and encoding will panic.
- All data integrity and validation will panic.

## Do one thing at a time and make sure that thing run very fast

- Things that falls into this category:
    - isolated decoding and encoding logic

## Concurrency

### Share Memory Model and Message Passing

- When it comes to non-trivial or complex synchronization, use message passing to 
isolate share access and ownership from thread instead of sharing memory model, or 
solely using very low level primitives.
- Regarless using message protocol or sharing memory model, think about using sane 
pattern instead of patterns that are designed for worst case scneario, and hope for 
the best. Example, don't soley rely on spanning mutex for every share data and hope 
they can guard against everything (It doesnt' work because: 1. guarding racing 
condition vs. guarding data integrity are completely two different things; 2. order 
of locking and deadlock, especially # of locks increase as there are more share data 
that need different level of thread safety; ...)
- Message protocol has X amount cost and overhead in different context. It isn't silver 
bullet. It's a treat off between performance and ease of implementing correct and safe 
synchronization. It's easy to fall in the trap of defaulting message protocol for every 
use case because of the lack of control and skill to make good use of mutex on share 
memory mdodel, and good understanding of the targeted problem.
- Message protocol and Share memory model are orthogonal but this doesn't mean they can't 
be use in hybrid.
- Favourite share memory model is trivial synchronization. Limit share memory model in 
use case where scope is fairly well defined.

### Mutex

- Mutex should have a clear intend, a clear usage of scope, a clear scope of data 
it needs to protect instead of using one single mutex to cover all things.

## Memory Management

- Create and destrory item are differ from allocating and freeing.
- Allocation should be visible to the caller. However, allocation detail should 
be hidden to the caller, and capture by the system.
- Transient object should be stack allocated. Transient object should not escape 
into the heap.
- If an object is long lived, heap allocate in the first place and make the 
allocation extremely obvious to locate and identify.
- Typical cache line size is 64 bytes

### Possible ways to make allocation to be visible for callers
        
- ...

## Coupling between `struct` and data stored in `struct`

- Regardless of using `struct` method or using regular procedures that accept 
a pointer of a `struct`, data should have minimal dependency with a specific 
`struct`. 
- Specifically, there some cases required for attention:
    - When combining multiple pieces of data to perform some work, there shouldn't 
    be case of which procedures needs to receive multiple different types of `struct`, 
    and use limited amount of field from each type of `struct`.
    - In order to access some pieces of data, there shouldn't be case of which 
    that pieces of data always need to get involve with a specific `struct`.
    - ...
- Procedure is a superset of functions and subroutines. A proceudre may or may 
not return something. A procedure may or may not have side effects.
    - A function is a mathematical entity that has no side effects.
    - A subroutine is something that has side effects but does not return 
    anything.
- Every procedure should follow the pattern of input -> process -> output. 
    - Ideal implemention should be input should also be the receiver of the output 
    for immutability. `a = process(a)`. This refers as (pure) function in Odin.
    - Less ideal situation is performing modification in place with pointer. Keep 
    in mind that this can be easily leaded back into building up coupling between 
    `struct` and data stored in `struct`. This refers as subroutine in Odin.
- The valid reason to tie `struct` and data stored in `struct` is getter and 
setter with non business logic validation (e.g. pointer not nil, duplication of 
key, mutex).

## Data Ownership and Data Read / Write Access

- Places that store data have the ownership.
- Read / write access (other than 1 to 1 simple assignment) should happen outside 
of places that store data.
- Read / write access outside of places that store data should deal of a view 
of data. That view typically can be data itself, or a combination of different 
data.
- Read / write access logic should be centralized and grouped, or / and 
encapsulated by a system (a set of procedures).
- Write access logic should produce output instead of performing modification in 
place.

## Data Correctness

### Encoding

#### Hierarchy

- Run a first pass on all hierarchies to make sure they all pass assertion test. 
Meanwhile, this first pass will also calculate size field in the HIRC's chunk header.
- Run a second pass to perform encoding. No size calculation and assertion test (
except size assertion test during the development phase to ensure encoding is bug 
free)
