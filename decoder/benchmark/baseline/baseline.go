package baseline

import (
	"bufio"
	"fmt"
	"io"
	"log"
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

func WriteOnce(data []byte) (f *os.File, err error) {
	f, err = os.CreateTemp("", "tmp")
	if err != nil {
		return nil, err
	}
	_, err = f.Write(data)
	return f, err
}

func UnBufferWrite(data []byte, writeSize uint32) (f *os.File, err error) {
	f, err = os.Create("tmp")
	if err != nil {
		return nil, err
	}

	for len(data) > 0 && err != nil {
		if int(writeSize) > len(data) {
			_, err = f.Write(data)
		} else {
			_, err = f.Write(data[:writeSize])
			data = data[writeSize:]
		}
	}

	return f, err
}

func BufferWrite(data []byte, bufSize uint32) (f *os.File, err error) {
	f, err = os.CreateTemp("", "tmp")
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriterSize(f, int(bufSize))

	nwrite := 0
	retry := 32
	for nwrite < len(data) && retry > 0 {
		nn, err := writer.Write(data)
		if err != nil {
			log.Print(err.Error())
			retry -= 1
		}
		nwrite += nn
	}

	if err = writer.Flush(); err != nil {
		return f, err
	}

	if nwrite < len(data) {
		return f, fmt.Errorf("Failed to write all data into the file")
	}

	return f, nil
}
