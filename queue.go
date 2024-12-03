package main

import (
	"log"
	"sync"
	"time"
)

var QueueMutex sync.Mutex
var Queue []int64

func RunQueue() {
	var program *Program
	{
		QueueMutex.Lock()
		defer QueueMutex.Unlock()

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
			return
		}
		// add it to the back.
		Queue = append(Queue, prog_id)
		if prog.disabled {
			return
		}
		log.Printf("running %v by %v", prog.name, prog.author)
		program = prog
	}

	ActiveProgram = program
	BroadcastEvent(InfoEvent())
	ExecuteProgram(program)
}

func ManageQueue() {
	for {
		RunQueue()
		// give it a break
		<-time.NewTimer(1 * time.Second).C
	}
}
