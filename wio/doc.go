package wio

// TYPE_SID: ShortID (uint32_t)
// TYPE_TID: target ShortID (uint32_t) same thing but for easier understanding of output
// TYPE_UNI: union (float / int32_t)
// TYPE_D64: double
// TYPE_F32: float
// TYPE_4CC: FourCC
// TYPE_S64: int64_t
// TYPE_U64: uint64_t
// TYPE_S32: int32_t
// TYPE_U32: uint32_t
// TYPE_S16: int16_t
// TYPE_U16: uint16_t
// TYPE_S8 : int8_t
// TYPE_U8 : uint8_t
// TYPE_VAR: variable size #u8/u16/u32
// TYPE_GAP: byte gap
// TYPE_STR: string
// TYPE_STZ: string (null-terminated)

// A helper struct that wraps an io.ReadSeeker. It provides a set of short 
// hand functions to read from this io.ReadSeeker, and produce commonly seen 
// data types.
// This helper struct is not designed for concurrent read with zero write and 
// zero copy. It's designed for an io.ReadSeeker that operates on bytes which 
// are not in the memory.
// If an io.ReadSeeker is operates on bytes that are in memory, such as 
// bytes.Reader, the following function will return copies instead of slices 
// that point to the same memory region:
// - Reader.ReadNUnsafe
// - Reader.ReadN, 
// - Reader.ReadAllUnsafe, 
// - Reader.ReadAll, 
