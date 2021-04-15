package pool

import "fmt"

type State int

const (
	StateOk State = iota + 1
	StateErr
)

func (s State) String() string {
	switch s {
	case StateOk:
		return "ok"
	case StateErr:
		return "err"
	default:
		return fmt.Sprintf("unknow state: %d", s)
	}
}
