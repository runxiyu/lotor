package main

import (
	bare "git.sr.ht/~runxiyu/go-bareish"
)

func decodeMsg(data []byte, val *Msg) error {
	return bare.Unmarshal(data, val)
}

func encodeMsg(val *Msg) ([]byte, error) {
	return bare.Marshal(val)
}

func encodeSend(send func([]byte) error, val interface{ bare.Union }) error {
	msg := Msg(val)
	return encodeMsgSend(send, &msg)
}

func encodeMsgSend(send func([]byte) error, val *Msg) error {
	data, err := encodeMsg(val)
	if err != nil {
		return err
	}
	err = send(data)
	if err != nil {
		return err
	}
	return nil
}
