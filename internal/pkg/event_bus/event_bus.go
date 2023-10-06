package event_bus

import (
	"errors"
	"eventsyncprocess/internal/pkg/event"
	"sync"
)

type eventChannel chan<- event.IEvent

func NewEventChannel(buffer int) chan event.IEvent {
	return make(chan event.IEvent, buffer)
}

type eventsMap map[string] /*eventName*/ eventChannel
type subscribersMap map[string] /*subscriberName*/ eventsMap

type EventBus struct {
	subs      subscribersMap
	subsMutex sync.Mutex
}

func NewEventBus() EventBus {
	return EventBus{
		subs:      make(subscribersMap),
		subsMutex: sync.Mutex{},
	}
}

//eventsCh := make(chan pubsub.Event, 3)
//assert.NoError(t, ps.Subscribe(sessionUUID, eventsCh))

func (b *EventBus) validateEventInput(eventName *string) error {
	if *eventName == "" {
		return errors.New("empty event name")
	}
	return nil
}

func (b *EventBus) validateSubscriberInput(subscriberName, eventName *string, channel chan<- event.IEvent) error {
	if *subscriberName == "" {
		return errors.New("empty subscriber name")
	}
	if err := b.validateEventInput(eventName); err != nil {
		return err
	}
	if channel == nil {
		return errors.New("event channel is nil")
	}
	return nil
}

func (b *EventBus) Subscribe(subscriberName, eventName string, channel chan<- event.IEvent) error {
	b.subsMutex.Lock()
	defer b.subsMutex.Unlock()

	err := b.validateSubscriberInput(&subscriberName, &eventName, channel)
	if err != nil {
		return err
	}

	if _, ok := b.subs[subscriberName]; !ok {
		b.subs[subscriberName] = make(eventsMap)
	}

	if _, ok := b.subs[subscriberName][eventName]; ok {
		return errors.New("already subscribed on event")
	}

	b.subs[subscriberName][eventName] = channel
	return nil
}

func (b *EventBus) Publish(eventName string, event event.IEvent) error {
	b.subsMutex.Lock()
	defer b.subsMutex.Unlock()

	//fmt.Println("Publish event : ", event.Name())
	err := b.validateEventInput(&eventName)
	if err != nil {
		return err
	}

	for _, subscriber := range b.subs {
		if targetCh, ok := subscriber[eventName]; ok {
			go func(c eventChannel) {
				c <- event
			}(targetCh)
		}
	}
	return nil
}
