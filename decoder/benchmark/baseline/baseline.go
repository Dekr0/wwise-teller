package baseline

import (
	"bufio"
	"io"
	"os"
)

func ReadOnce(p string) {
	data, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	for i := range data {
		data[i] = 0
	}
}

func BufferRead(recvSize uint32, bufSize uint32, p string) {
	buf := make([]byte, recvSize, recvSize)

	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, int(bufSize))

	for {
		_, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}

func UnbufferRead(recvSize uint32, p string) {
	buf := make([]byte, recvSize, recvSize)

	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		_, err = f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}
