package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sort"
	"time"
)

type Process struct {
	id      int
	current int
	active  bool
}

func (p *Process) execute() {
	i := 0
	for {
		if p.active == false {
			break
		}
		p.current++
		fmt.Println("\nID:", p.id, "-", p.current)
		time.Sleep(time.Millisecond * 500)
		i++
	}
}

func getProcess(processes []*Process) *Process {
	sort.Slice(processes[:], func(i, j int) bool {
		return processes[i].id < processes[j].id
	})
	p := processes[0]
	copy(processes[0:], processes[1:])
	processes = processes[:len(processes)-1]
	processes = processes[0:]
	p.active = false
	return p
}

func server(processes []*Process) {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(processes) > 0 {
			go handleClient(conn, processes)
		}
	}
}

func handleClient(conn net.Conn, processes []*Process) {
	sort.Slice(processes[:], func(i, j int) bool {
		return processes[i].id < processes[j].id
	})
	p := processes[0]
	copy(processes[0:], processes[1:])
	processes = processes[:len(processes)-1]
	processes = processes[0:]
	p.active = false
	err := gob.NewEncoder(conn).Encode(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	var decodedProcess Process
	err = gob.NewDecoder(conn).Decode(&decodedProcess)
	if err != nil {
		fmt.Println(err)
		return
	}
	if &decodedProcess != nil {
		decodedProcess.active = true
		go decodedProcess.execute()
		processes = append(processes, &decodedProcess)
		sort.Slice(processes[:], func(i, j int) bool {
			return processes[i].id < processes[j].id
		})
	}
}

func initProcesses(n int) []*Process {
	var processes []*Process
	for i := 0; i < n; i++ {
		p := Process{id: i, current: 0, active: true}
		go p.execute()
		processes = append(processes, &p)
	}
	return processes
}

func main() {
	go server(initProcesses(5))

	var input string
	fmt.Scanln(&input)
}
