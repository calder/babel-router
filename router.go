package main

import "errors"

type Router struct {
    queue MsgQueue
    sendQueues map[string]MsgQueue
}

type MsgQueue chan []byte

func NewRouter () *Router {
    router := &Router{make(MsgQueue), make(map[string]MsgQueue)}
    go router.processQueue()
    return router
}

func newTestRouter () *Router {
    return &Router{make(MsgQueue), make(map[string]MsgQueue)}
}

func (r *Router) enqueue (msg []byte) {
    r.queue <- msg
}

func (r *Router) hasSendQueue (s string) bool {
    _, ok := r.sendQueues[s]
    return ok
}

func (r *Router) openSendQueue (s string) MsgQueue {
    if _, ok := r.sendQueues[s]; ok {
        log.Debug("send queue already open for %s", SHex(s))
        return nil
    }

    q := make(MsgQueue)
    r.sendQueues[s] = q
    return q
}

func (r *Router) closeSendQueue (s string) {
    delete(r.sendQueues, s)
}

func (r *Router) addToSendQueue (dest string, msg []byte) error {
    if q, ok := r.sendQueues[dest]; !ok {
        return errors.New("destination not connected: " + SHex(dest))
    } else {
        q <- msg
        return nil
    }
}
