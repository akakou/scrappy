package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	"github.com/akakou/ecdaa/tpm_utils"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/ecdaa_helper"
	_ "github.com/mattn/go-sqlite3"
)

const PASSWORD = "password"

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

func writeError(err error) {
	var resp Response

	resp.Status = "error"
	resp.Error = err.Error()
	write(&resp)
}

func main() {
	var resp Response

	req, err := read()

	if err != nil {
		writeError(err)
		return
	}

	tpm, err := tpm_utils.OpenTPM([]byte(PASSWORD), "/dev/tpm0")
	if err != nil {
		writeError(err)
		return
	}

	ecdaa_helper.SIGNER_TPM_CONFIG_PATH = "/scrappy/signer-tpm.json"
	signer, err := ecdaa_helper.PrepareTPMSignerFromFile(tpm)
	if err != nil {
		writeError(err)
		return
	}

	signature, err := scrappy.Sign(req.Origin, req.Period, signer)

	if err != nil {
		writeError(err)
		return
	}

	resp.Status = "ok"
	resp.Signature = signature

	write(&resp)
}
