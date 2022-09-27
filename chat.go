package main

import (
	"bufio"
	"fmt"
	"github.com/GinGin3203/protohackers/pkg/must"
	"io"
	"log"
	"net"
	"time"
	"unicode"
)

type Manager struct {
}

func newManager() *Manager {
	return &Manager{}
}

func (m *Manager) Serve(c net.Conn) {
	defer c.Close()
	must.NoError(c.SetDeadline(time.Now().Add(10 * time.Second)))
	userName, err := namePrompt(c)
	if err != nil {
		log.Println(err)
		c.Write([]byte(err.Error()))
		return
	}
	m.AnnounceNewUser(userName)
}

func (m *Manager) AnnounceNewUser(name string) {

}

func namePrompt(c net.Conn) (string, error) {
	must.NoError(c.SetDeadline(time.Now().Add(10 * time.Second)))
	if _, err := c.Write([]byte(fmt.Sprintln("Enter your name: "))); err != nil {
		log.Println(err)
		return "", err
	}
	rd := bufio.NewReader(io.LimitReader(c, 64))
	b, err := rd.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return "", err
	}
	if !isASCII(b) {
		err = fmt.Errorf("non-ascii name: %s", string(b))
		return "", err
	}
	return string(b[:len(b)-1]), nil
}

func isASCII(s []byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
