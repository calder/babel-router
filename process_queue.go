package main

import "encoding/hex"
import "github.com/calder/babel-lib-go"

// Message types
var BLOB = babel.Type("01")

type Context struct {
    dest *babel.Id1
}

func (r *Router) handleBlob (data []byte, c *Context) {
    log.Debug("blob: %s", hex.EncodeToString(data))
    if c.dest == nil { log.Debug("message has no recipient"); return }
}

func (r *Router) handleEnvelope (data []byte, c *Context) {
    env, e := DecodeEnvelope(data)
    if e != nil { log.Debug("error decoding envelope: %s", e); return }
    destType, destData, _, e := babel.Unwrap(babel.TYPE, env.dest)
    if destType != babel.ID1 { log.Debug("unknown destination type: %x", destType); return }
    dest, e := babel.DecodeId1(destData)
    if e != nil { log.Debug("error decoding destination: %s", e); return }
    c.dest = dest
    r.handleMessage(env.msg, c)
}

func (r *Router) handleMessage (msg []byte, c *Context) {
    typ, data, _, e := babel.Unwrap(babel.TYPE, msg)
    if e != nil { log.Debug("error decoding message: %s", e); return }
    switch typ {
    case BLOB:     r.handleBlob(data, c)
    case ENVELOPE: r.handleEnvelope(data, c)
    default:       log.Debug("unknown message type: %d", typ)
    }
}

func (r *Router) processQueue () {
    for {
        select {
        case msg := <-r.queue:
            if msg == nil { break }
            r.handleMessage(msg, &Context{})
        }
    }
}
