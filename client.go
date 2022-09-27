package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

type Listener struct {
	name username
	net.Conn
	newUser        <-chan userJoinedAnnouncement
	sendMessage    chan<- message
	receiveMessage <-chan message
	done           chan struct{}
}

func NewClient(u userJoinedAnnouncement) *Listener {
	return &Listener{
		name:    u,
		newUser: make(<-chan userJoinedAnnouncement),
	}
}

func (c *Listener) Start() {
	go c.ReadAsync()
	for {
		select {
		case un := <-c.newUser:
			if _, err := c.Write(un); err != nil {
				log.Print(err)
				return
			}
		case msg := <-c.receiveMessage:
			if _, err := c.Write(msg); err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func (c *Listener) ReadAsync() {
	sc := bufio.NewScanner(&io.LimitedReader{R: c, N: 1024})
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		tag := make([]byte, len(c.name)+3)
		tag[0] = '['
		copy(tag[1:len(c.name)+1], c.name)
		tag[len(tag)-2] = ']'
		tag[len(tag)-1] = ' '
		c.sendMessage <- append(tag, sc.Bytes()...)
	}
	log.Println(sc.Err())
	c.done <- struct{}{}
}
