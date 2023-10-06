package first_step

import "eventsyncprocess/step"

type firstStep struct {
	Name string
}

func NewStep() step.IStep {
	return &firstStep{}
}

func (s *firstStep) Execute() {

}
