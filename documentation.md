# Keybinding

- I will use `mod` as an alias for `ctrl+shift`.

## Navigation

### File Explorer and File Dialog

- `ctrl+f` - focus on the search bar
- `mod+f` - focus on the first entry
- `j` - go to next (down) entry
- `k` - go to previous (up) entry
- `h` - go to the parent directory
- `k` - go to the selected directory, or open the file (for file explorer only)
- `shift+j` - multiple select down
- `shift+k` - multiple select up
- `ctrl+a` - select all
- `ctrl+s` - open the selected entries
- `ctrl+q` - close the file dialog

### Command Palette

- `ctrl+f` - focus on the search bar
- `mod+f` - focus on the first entry
- `j` - go to next (down) entry
- `k` - go to previous (up) entry
- `ctrl+s` - execute the selected command entry 
- `ctrl+q` - close the command palette

### Docking Windows

- `mod+j` - cycle to the previous docking window
- `mod+k` - cycle to the next docking window
- `ctrl+tab` - open docking window tab, and cycle through each docking window

### Bank Explorer

#### Navigation

- `ctrl+f` - focus on the search bar
- `mod+f` - focus on the first entry
- `j` - go to next (down) entry
- `k` - go to previous (up) entry
- `shift+j` - multiple select down
- `shift+k` - multiple select up
- `ctrl+a` - select all

#### File Operation

- `mod+s` - save the sound bank that is currently displayed in the bank explorer 
without any type of integration.
- `mod+i` - generate a Helldivers 2 patch that contains the sound bank that is 
currently displayed in the bank explorer

# Notes

## Values Bound

- There are lower bound and upper bound for different values despite the input 
fields for floating point value / integer do not validate this due to lack of 
documentation from Wwise.
- For example, make up gain can only take value in between -96.0 and 12.0
- Some values are hierarchy ID based. Currently, I haven't implemented the 
ability of inputting these types of values by drop down selection because it 
requires me to scan through all sound banks.
    - Attenuation ID is one example. The value is sane only when this value is 
    associated with an known Attenuation hierarchy object. This attenuation 
    hierarchy object can be in the same sound bank, or it can be in a complete 
    different sound bank.

## Values Type of Property and Range Property 

- Some properties will only make use floating point type (`f32`) while others 
make use of integer type, which can be signed (`int32`) or unsigned (`uint32`).
- The correct type of each property is not documented. Feel free to document 
this by checking each property in Wwise Authoring tool.

### Listing

| Property          | Type   | Min   | Max  | Notes                                                        |
|-------------------|--------|-------|------|--------------------------------------------------------------|
| Make Up Gain      | f32    | -96.0 | 12.0 | n/a                                                          |
| Output Bus Volume | f32    | -96.0 | 12.0 | n/a                                                          |
| Initial Delay     | f32    | ?     | ?    | n/a                                                          |
| Attenuation ID    | uint32 | n/a   | n/a  | sane only if an attenuation hierarchy object has the same ID |
