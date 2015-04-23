package main

import "encoding/hex"
import "github.com/calder/babel-lib-go"

var ENVELOPE = babel.Type("EBB7B521")
func (*Envelope) Type () uint64 { return ENVELOPE }
func (*Envelope) TypeName () string { return "Envelope" }
func init () { babel.AddType(ENVELOPE, decodeEnvelope) }

type Envelope struct {
    dest []byte
    msg  []byte
}

func NewEnvelope (dest []byte, msg []byte) *Envelope {
    return &Envelope{dest, msg}
}

func EnvelopeFromValues (dest babel.Value, msg babel.Value) *Envelope {
    return NewEnvelope(dest.CBR(), msg.CBR())
}

func (e *Envelope) String () string {
    return "<Envelope:" + hex.EncodeToString(e.dest) + "," + hex.EncodeToString(e.msg) + ">"
}

func (e *Envelope) CBR () []byte {
    return e.Encode(babel.TYPE)
}

func (e *Envelope) Encode (enc babel.Encoding) []byte {
    return babel.Wrap(enc, ENVELOPE, babel.Join(babel.Wrap(babel.LEN, 0, e.dest), e.msg))
}

func decodeEnvelope (data []byte) (babel.Value, error) { return DecodeEnvelope(data) }
func DecodeEnvelope (data []byte) (*Envelope, error) {
    _, dest, n, e := babel.Unwrap(babel.LEN, data)
    if e != nil { return nil, e }

    msg := data[n:]
    return &Envelope{dest, msg}, nil
}
