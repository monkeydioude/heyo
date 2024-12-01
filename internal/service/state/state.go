package state

type State int

const (
	STATE_IDLE State = iota
	STATE_BUSY       = iota
)

func (s State) Idle() {
	s = STATE_IDLE
}

func (s State) Busy() {
	s = STATE_BUSY
}

func Idle() State {
	var s State
	s.Idle()
	return s
}

func Busy() State {
	var s State
	s.Busy()
	return s
}
