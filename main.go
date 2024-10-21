package main

import (
	"bufio"
	"os"
)

func main() {
	// Temporary sender since we don't have networking code yet
	writer := bufio.NewWriter(os.Stdout)
	send := func(data []byte) error {
		_, err := writer.Write(data)
		if err != nil {
			return err
		}
		err = writer.Flush()
		return err
	}

	msg := MsgJoin{
		Target: []byte("target"),
	}
	err := encodeSend(send, msg)
	if err != nil {
		panic(err)
	}
}
