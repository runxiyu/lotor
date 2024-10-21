package main

import "git.sr.ht/~sircmpwn/go-bare"

func decodeMsg(data []byte, val *Msg) error {
	return bare.Unmarshal(data, val)
}

func encodeMsg(val *Msg) ([]byte, error) {
	return bare.Marshal(val)
}
