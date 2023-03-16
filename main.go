package main

import (
	"fmt"
)

type Person interface {
	echo()
	set(string)
}

type p1 struct {
	name string
}

func (p *p1) echo() {
	fmt.Println(p.name)
}

func (p *p1) set(name string) {
	p.name = name
}

type p2 struct {
	Person
}

func main() {
	tian := &p2{&p1{"ye"}}
	tian.echo()
}
