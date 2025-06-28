# Database

## Source IDing

- Source IDs should be persistent if audio sources don't change at all. This 
reduce collision of source IDs. This is similar to how Wwise keeps track of 
source ID in a project.

## Hierarchy Object IDing

- Hierarchy object IDs should be persistent. This persistence can be determined 
by users label.

# Reactive UI System

# Draft Idea

- Fine grained synchronous reactivity using reactive primitive such as signal
    - Callback chainging
- Event bus / Event streaming asynchronous reactivity
    - producer emit data into buses / streams
    - emit data into buses / streams
    - buses / streams transform into new types of buses / streams by applying 
    functional like primitives such as `map`, `buffer`, `reduce`, etc.
    - consumer consume data emitted from buses / streams
