package wio

import "core:io"
import "core:bufio"

wio_assert :: proc()

bufio_reader_must_u8 :: proc(b: ^bufio.Reader, big: bool) -> u8 {
    data, err := bufio.reader_read_byte(b)
    if err != 0 {}
}
