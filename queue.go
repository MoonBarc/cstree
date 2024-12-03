package main

import (
	"log"
	"sync"
	"time"
)

var QueueMutex sync.Mutex
var Queue []int64 = make([]int64, 0)

func RunQueue() {
	var program *Program
	QueueMutex.Lock()

	if len(Queue) == 0 {
		return
	}

	// pop the next one
	prog_id := Queue[0]
	Queue = Queue[1:]
	// fetch code, if disabled
	prog, err := TrDB.GetProgram(prog_id)
	if err != nil {
		log.Println("warning, failed to get prog of id", prog.id, err)
		QueueMutex.Unlock()
		return
	}
	// add it to the back.
	Queue = append(Queue, prog_id)
	QueueMutex.Unlock()
	if prog.disabled {
		return
	}
	log.Printf("running %v by %v", prog.name, prog.author)
	program = prog

	ActiveProgram = program
	BroadcastEvent(InfoEvent())
	ExecuteProgram(program)
}

func ManageQueue() {
	progs, err := TrDB.AllPrograms()
	if err != nil {
		log.Fatalln("failed to load programs into queue", err)
	}
	for _, p := range progs {
		Queue = append(Queue, p.id)
	}
	for {
		RunQueue()
		// give it a break
		<-time.NewTimer(1 * time.Second).C
	}
}
