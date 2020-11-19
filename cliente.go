package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Process struct {
	id      int
	current int
	active  bool
}

func (p *Process) client() {
	i := 0
	for {
		p.current++
		fmt.Println("ID:", p.id, "-", p.current)
		time.Sleep(time.Millisecond * 500)
		i++
	}
}

func main() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var p Process
	err = gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	go p.client()

	var input string
	fmt.Scanln(&input)

	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}
