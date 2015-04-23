package main

import "bytes"
import "encoding/hex"
import "github.com/calder/babel-lib-go"

var BLOB = babel.Type("01")
func (*Blob) Type () uint64 { return BLOB }
func (*Blob) TypeName () string { return "Blob" }
func init () { babel.AddType(BLOB, decodeBlob) }

type Blob struct {
    data []byte
}

func NewBlob (data []byte) *Blob {
    return &Blob{data}
}

func BlobFromValue (content babel.Value) *Blob {
    return NewBlob(content.CBR())
}

func (b *Blob) Data () []byte {
  return b.data
}

func (b *Blob) String () string {
    return "<Blob:" + hex.EncodeToString(b.data) + ">"
}

func (b *Blob) CBR () []byte {
    return b.Encode(babel.TYPE)
}

func (b *Blob) Encode (enc babel.Encoding) []byte {
    return babel.Wrap(enc, BLOB, b.data)
}

func decodeBlob (data []byte) (babel.Value, error) { return DecodeBlob(data), nil }
func DecodeBlob (data []byte) *Blob {
    return &Blob{data}
}

func (b *Blob) Equal (other *Blob) bool {
    return bytes.Equal(b.data, other.data)
}
