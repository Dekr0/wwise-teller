# Wwise Teller

## Summary

- Wwise teller is a (modding) SDK for editing encoded sound bank file generated 
from Wwise Authoring Tool. 
- Wwise teller makes an attempt to replicate common features sound designer can 
do in Wwise Authoring Tool
- A quick [sneak peek](https://youtu.be/36MphHqG2ks](https://youtu.be/36MphHqG2ks)) on the current state of 
Wwise Teller

## Limitation

- Wwise Teller is still at its very eariler stage of development. The current 
goal of Wwise Teller is to have the abilities to edit encoded sound bank file 
with version 141.
- Games that make use of sound bank version 141:
  - Helldivers 2,
  - Overwatch ?
  - ...

## Usage

- Detail check [here](usage.md)

## Contribution

- Feel free to leave suggestions or PRs (especially on performance and state 
management / pattern in the UI part)
- Wwise Teller might undergo a direct port to Zig or Odin if it needs to start 
doing any sort operations that performance intensive such as real time audio 
processing, audio simulation etc.

### Code Organization

- `assert` - hand roll assertion function
- `interp` - mathematics interpolation for things such as visualizing RTPC, 
Modulator, etc.
- `parser` - Sound bank parser
- `ui` - UI logic
    - files with prefix `re` are render related.
    - files with prefix `st` are `structs` that contain state used in the render.
    - files with prefix `cb` are callback creation for different UI widgets
- `wio` - `structs` that make IO and encoding / decoding easier
- `wwise` - `structs` for data in Sound banks
