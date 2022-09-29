package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"
)

type Connection struct {
	net.Conn
	name []byte
	ch   chan event
}

func (c *Connection) ListenFromServer(ctx context.Context) error {
	for {
		select {
		case msg := <-c.ch:
			var err error
			err = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				return err
			}
			if _, err = c.Write(msg.content); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Connection) SendToServer() error {
	sc := bufio.NewScanner(&io.LimitedReader{R: c, N: 1024})
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		c.ch <- newMessage(append(tag(c.name), sc.Bytes()...))
	}
	return sc.Err()
}

func tag(name []byte) []byte {
	tag := make([]byte, len(name)+3)
	tag[0] = '['
	copy(name, tag[1:len(name)+1])
	tag[len(tag)-2] = ']'
	tag[len(tag)-1] = ' '
	return tag
}
