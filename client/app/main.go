package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Request struct {
	Time   string `json:"time"`
	Origin string `json:"origin"`
}

type Response struct {
	Status    string `json:"status"`
	Signature string `json:"signature"`
	Error     string `json:"error"`
}

func main() {
	var inLength uint32
	var req Request
	var resp Response

	stdin := bufio.NewReader(os.Stdin)
	binary.Read(stdin, binary.LittleEndian, &inLength)
	buf := make([]byte, inLength)
	io.ReadFull(stdin, buf)

	err := json.Unmarshal(buf, &req)

	if err != nil {
		resp = Response{
			Status:    "error",
			Signature: "",
			Error:     err.Error(),
		}
	} else {
		resp = Response{
			Status:    "ok",
			Signature: fmt.Sprintf("signature for %s", req.Time),
			Error:     "",
		}
	}

	payload, _ := json.Marshal(&resp)

	stdout := bufio.NewWriter(os.Stdout)
	outLength := len(payload)

	binary.Write(stdout, binary.LittleEndian, int32(outLength))
	for head := 0; head < outLength; {
		n, _ := stdout.Write(payload[head:])
		head += n
	}

	stdout.Flush()
}
