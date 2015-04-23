package main

import "github.com/calder/babel-lib-go"

type Message struct {
    dest *babel.Hash1
}

func (r *Router) handleBlob (data []byte, m *Message) {
    blob, e := DecodeBlob(data)
    if e != nil { log.Debug("error decoding blob: %s", e); return }
    log.Debug("got %s", blob)

    if m.dest == nil { log.Debug("message has no recipient"); return }

    e = r.addToSendQueue(string(m.dest.CBR()), blob.Data())
    if e != nil { log.Debug("error sending message: %s", e); return }
}

func (r *Router) handleEnvelope (data []byte, m *Message) {
    env, e := DecodeEnvelope(data)
    if e != nil { log.Debug("error decoding envelope: %s", e); return }
    log.Debug("got %s", env)

    destType, destData, _, e := babel.Unwrap(babel.TYPE, env.dest)
    if destType != babel.ID1 { log.Debug("unknown destination type: %x", destType); return }

    dest, e := babel.DecodeHash1(destData)
    if e != nil { log.Debug("error decoding destination: %s", e); return }

    m.dest = dest
    r.handleMessage(env.msg, m)
}

func (r *Router) handleMessage (msg []byte, m *Message) {
    typ, data, _, e := babel.Unwrap(babel.TYPE, msg)
    if e != nil { log.Debug("error decoding message: %s", e); return }
    switch typ {
    case BLOB:     r.handleBlob(data, m)
    case ENVELOPE: r.handleEnvelope(data, m)
    default:       log.Debug("unknown message type: %d", typ)
    }
}

func (r *Router) processQueue () {
    for {
        select {
        case msg := <-r.queue:
            if msg == nil { break }
            r.handleMessage(msg, &Message{})
        }
    }
}
