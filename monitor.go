package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Event struct {
	Ty   string
	Data map[string]any
}

type Conn struct {
	Rx chan Event
}

var Conns map[string]*Conn = make(map[string]*Conn)
var ConnsMutex sync.Mutex

func BroadcastEvent(ev Event) {
	ConnsMutex.Lock()
	defer ConnsMutex.Unlock()
	for _, conn := range Conns {
		conn.Rx <- ev
	}
}

func AddConnection(remoteAddr string, conn *Conn) {
	ConnsMutex.Lock()
	defer ConnsMutex.Unlock()
	Conns[remoteAddr] = conn
}

func RemoveConnection(remoteAddr string) {
	ConnsMutex.Lock()
	defer ConnsMutex.Unlock()
	delete(Conns, remoteAddr)
}

func InfoEvent() Event {
	p := ActiveProgram
	return Event{
		Ty: "info",
		Data: map[string]any{
			"title":  p.name,
			"author": p.author,
		},
	}
}

func (c *Conn) SendInfo() {
	c.Rx <- InfoEvent()
}

func Monitor(w http.ResponseWriter, r *http.Request) {
	conn := Conn{
		Rx: make(chan Event, 10),
	}

	AddConnection(r.RemoteAddr, &conn)
	defer RemoveConnection(r.RemoteAddr)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	conn.SendInfo()

	// Flush allows streaming responses
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	finished := r.Context().Done()

	for {
		select {
		case event := <-conn.Rx:
			// JSON encode the event
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error encoding event: %v", err)
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

		case <-finished:
			log.Println("client disconnected")
			return
		}
	}
}
