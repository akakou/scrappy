package main

import (
	"bufio"
	"core"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Request struct {
	Period int    `json:"period"`
	Origin string `json:"origin"`
}

type Response struct {
	Status    string `json:"status"`
	Signature string `json:"signature"`
	Error     string `json:"error"`
}

func read() (*Request, error) {
	var inLength uint32
	var req Request

	stdin := bufio.NewReader(os.Stdin)
	binary.Read(stdin, binary.LittleEndian, &inLength)
	buf := make([]byte, inLength)
	io.ReadFull(stdin, buf)

	err := json.Unmarshal(buf, &req)

	return &req, err
}

func write(resp *Response) error {
	payload, _ := json.Marshal(&resp)

	stdout := bufio.NewWriter(os.Stdout)
	outLength := len(payload)

	binary.Write(stdout, binary.LittleEndian, int32(outLength))
	for head := 0; head < outLength; {
		n, _ := stdout.Write(payload[head:])
		head += n
	}

	stdout.Flush()

	return nil
}

func main() {
	var resp Response

	req, err := read()

	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
		write(&resp)
		return
	}

	core.CONFIG_PATH = "/attestation/config.json"

	signature, err := core.Sign(req.Origin, req.Period)

	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
		write(&resp)
		return
	}

	resp.Status = "ok"
	resp.Signature = signature

	write(&resp)
}
