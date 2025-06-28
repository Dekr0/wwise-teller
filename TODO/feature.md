# Information Copying

- Every single ID

# File Dialog and File Explorer

- Modify the destination file name
- Ability to delete / rename file
- Split a file path into individual segment. When clicking on a segment, it will 
jump to a specific file location.

# Helldivers 2 Integration

- Pack multiple sound banks for Helldivers 2 integration

# Rewiring

## Change parent of a hierarchy & Add / Remove child from a container

- Update: this only works for node with no leafs!

### Change parent of a hierarchy

- Flow:
    1. The parent remove the target hierarchy.
    2. If step 1 success, the new parent add the target hierarchy.
    3. If step 2 success, change direct parent ID of the target hierarchy.
    4. Arrange the tree index. Put the target hierarchy next to the new parent

### Remove a child from a container

- Flow:
    1. The parent remove the target hierarchy.
    2. If step 1 success, change direct parent ID of the target hierarchy to 
    nothing.
    3. Arrange the tree index. Put the target hierarchy to free float.
