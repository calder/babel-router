package main

import "log"
import "github.com/calder/babel-lib-go"

func (r *Router) processMessage (msg []byte) {
    m, e := babel.Decode(msg)
    if e != nil { log.Println("error decoding message:", e); return }
    println(m)
}

func (r *Router) processQueue () {
    for {
        select {
        case msg := <-r.queue:
            if msg == nil { break }
            r.processMessage(msg)
        }
    }
}