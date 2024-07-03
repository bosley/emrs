package core

import (
	"fmt"
)

type Err struct {
	data *Stack[string]
}

func NErr(what string) *Err {
	return &Err{
		data: NewStack[string](""),
	}
}

func ErrFrom(err error) *Err {
	return NErr(err.Error())
}

func (e *Err) Push(what string) *Err {
	e.data.Push(what)
	return e
}

func (e *Err) Error() string {
	x := 0
	result := ""
	for {
		w, e := e.data.Pop()
		if e != nil {
			return result
		}
		result += fmt.Sprintf("[%d]: %s\n", x, w)
		x += 1
	}
	return ""
}
