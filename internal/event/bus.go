package event

type Event struct {
	Type string
	Data interface{}
}

var Bus = make(chan Event, 100) // buffered for non-blocking async
