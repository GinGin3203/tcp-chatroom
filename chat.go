package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/GinGin3203/protohackers/pkg/must"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"unicode"
)

type Manager struct {
	*sync.RWMutex
	pool map[chan event]struct{}
}

func newManager() *Manager {
	return &Manager{
		RWMutex: &sync.RWMutex{},
		pool:    map[chan event]struct{}{},
	}
}

func (m *Manager) Serve(c net.Conn) {
	userName, err := namePrompt(c)
	if err != nil {
		log.Println(err)
		c.Write([]byte(err.Error()))
		c.Close()
		return
	}
	connCh := make(chan event, 16)
	m.AddConnChannel(connCh)

	conn := &Connection{
		name: userName,
		Conn: c,
		ch:   connCh,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go conn.ListenFromServer(ctx)
	go conn.SendToServer()

	for msg := range connCh {
		m.RLock()
		for otherConnCh, _ := range m.pool {
			if connCh == otherConnCh {
				continue
			}
			otherConnCh <- msg
		}
		m.Unlock()
	}

	cancel()
}

func (m *Manager) Disconnect(c *Connection) {
	_ = c.Close()

	m.Lock()
	delete(m.pool, c.ch)
	m.Unlock()

	close(c.ch)
}

func (m *Manager) AddConnChannel(c chan event) {
	m.Lock()
	m.pool[c] = struct{}{}
	m.Unlock()
}

func (m *Manager) Broadcast(ctx context.Context, writer chan event) {
	for {
		select {
		case msg := <-writer:
			m.RLock()
			for s, _ := range m.pool {
				if s == writer {
					continue
				}
				s <- msg
			}
			m.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}

func namePrompt(c net.Conn) ([]byte, error) {
	must.NoError(c.SetDeadline(time.Now().Add(10 * time.Second)))
	if _, err := c.Write([]byte(fmt.Sprintln("Enter your name: "))); err != nil {
		log.Println(err)
		return nil, err
	}
	rd := bufio.NewReader(io.LimitReader(c, 64))
	b, err := rd.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !isASCII(b) {
		err = fmt.Errorf("non-ascii name: %s", string(b))
		return nil, err
	}
	return b, nil
}

func isASCII(s []byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
