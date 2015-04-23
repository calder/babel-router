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

func (r *Router) hasSendQueue (dest []byte) bool {
    _, ok := r.sendQueues[string(dest)]
    return ok
}

func (r *Router) openSendQueue (dest []byte) MsgQueue {
    if _, ok := r.sendQueues[string(dest)]; ok {
        log.Debug("send queue already open for %s", Hex(dest))
        return nil
    }

    q := make(MsgQueue)
    r.sendQueues[string(dest)] = q
    return q
}

func (r *Router) closeSendQueue (dest []byte) {
    delete(r.sendQueues, string(dest))
}

func (r *Router) addToSendQueue (dest []byte, msg []byte) error {
    if q, ok := r.sendQueues[string(dest)]; !ok {
        return errors.New("destination not connected: " + Hex(dest))
    } else {
        q <- msg
        return nil
    }
}
