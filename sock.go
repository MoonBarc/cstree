package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
)

// 2 = magic number
// 4*150 = 4 32bit colors
const packetLen = 2 + 4*150

const sockAddr = "/tmp/cstree.sock"

const magic = 0x1999

func DealWithSocket(c net.Conn) {
	log.Println("connection on socket!")
	defer c.Close()

OuterLoop:
	for {
		buf := make([]byte, packetLen)
		bytesRead := 0
		for bytesRead < packetLen {
			n, err := c.Read(buf[bytesRead:])
			if err != nil {
				log.Println("socket read err", err)
				break OuterLoop
			}
			bytesRead += n
		}
		attempted_magic := binary.LittleEndian.Uint16(buf[0:2])
		if attempted_magic != magic {
			log.Println("magic wrong, got", attempted_magic)
			break
		}

		// got a command!
		BroadcastEvent(Event{
			Ty: "render",
			Data: map[string]any{
				"colors": buf[2:],
			},
		})
		SetLEDs(buf[2:])
	}
	log.Println("disconnect from socket")
}

func ListenSocket() {
	_ = os.Remove(sockAddr)
	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatalln("failed to listen on socket", err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatalln("couldn't listen to socket??", err)
		}
		go DealWithSocket(fd)
	}
}
