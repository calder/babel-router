package main

import "github.com/calder/babel-lib-go"

var ENVELOPE = babel.Type("EBB7B521")

type Envelope struct {
    dest []byte
    msg  []byte
}

func DecodeEnvelope (data []byte) (*Envelope, error) {
    _, dest, n, e := babel.Unwrap(babel.LEN, data)
    if e != nil { return nil, e }
    msg := data[n:]
    return &Envelope{dest, msg}, nil
}
