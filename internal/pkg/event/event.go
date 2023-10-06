package event

type IEvent interface {
	Name() string
	Payload() interface{}
}

type event struct {
	name    string
	payload interface{}
}

func (e *event) Name() string         { return e.name }
func (e *event) Payload() interface{} { return e.payload }

func NewEvent(name string, payload interface{}) IEvent {
	return &event{
		name:    name,
		payload: payload,
	}
}
