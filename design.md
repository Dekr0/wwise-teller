# Design

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
