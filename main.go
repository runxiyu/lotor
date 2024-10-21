package main

import (
	"bufio"
	"os"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	cmdMsg := MsgJoin{
		Target: []byte("target"),
	}
	msg := Msg(cmdMsg)
	encoded, err := encodeMsg(&msg)
	if err != nil {
		panic(err)
	}
	_, err = writer.Write(encoded)
	if err != nil {
		panic(err)
	}
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
