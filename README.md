# Wwise Teller

## Credit

- Wwise Teller cannot be made without the helps of `wwiser`.

## Summary

- Wwise teller is a (modding) SDK for editing encoded sound bank file generated 
from Wwise Authoring Tool. 
- Wwise teller makes an attempt to replicate common features sound designer can 
do in Wwise Authoring Tool

## Capabilities

- The current capability of Wwise Teller (Most are done by scripting / automation at the current state) include the following:
    - Add new Audio Sources with new audio data (with completely new Audio Source 32-bit IDs)
    - Replace Audio Sources Data
    - Swap Audio source IDs in Sound objects
    - Add new Actions with specified targets, and append them to a given Event (Not fully complete)
    - Add any types of hierarchy in the Actor Mixer Hierarchy Categories, and wire them up with a
    new Action (Not fully complete)
    - Modify different properties commonly seen in Wwise Authoring Tool (Not fully complete)
        - Volume, Make Up Gain, Initial Delay, ...
        - User-Defined Auxiliary Send (Send Volume and the Auxiliary Send being used)
        - Playback Limit
        - Virtual Behavior
        - ...
    - Modify RTPC Curve
    - Modify Bus (Not fully complete)
    - Modify Attenuation (Not fully complete)
    - ...
- Notice that the above is proven to be working in Helldivers 2.

## Limitation

- Wwise Teller is still at its very eariler stage of development. The current 
goal of Wwise Teller is to have the abilities to edit encoded sound bank file 
with version 141 and version 154.
- Games that make use of sound bank version 141:
  - Helldivers 2,
  - Overwatch ?
  - ...

## Usage

- Please read the [wiki](https://github.com/Dekr0/wwise-teller/wiki)

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
