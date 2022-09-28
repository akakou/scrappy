package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

func main() {
	var inLength uint32

	stdin := bufio.NewReader(os.Stdin)
	binary.Read(stdin, binary.LittleEndian, &inLength)
	buf := make([]byte, inLength)
	io.ReadFull(stdin, buf)

	payload := []byte("{\"signature\":\"signature\",\"count\":0}")

	stdout := bufio.NewWriter(os.Stdout)
	outLength := len(payload)

	binary.Write(stdout, binary.LittleEndian, int32(outLength))
	for head := 0; head < outLength; {
		n, _ := stdout.Write(payload[head:])
		head += n
	}

	stdout.Flush()
}
