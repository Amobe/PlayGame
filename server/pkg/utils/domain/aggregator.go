package domain

type Aggregator interface {
	ID() string
	Events() []Event
	embedCoreAggregator()
}

type CoreAggregator struct {
	events []Event
}

func (a *CoreAggregator) Apply(events ...Event) {
	a.events = append(a.events, events...)
}

func (a *CoreAggregator) Events() []Event {
	res := a.events
	a.events = nil
	return res
}

func (a *CoreAggregator) embedCoreAggregator() {}

type Event interface {
	Name() string
	embedCoreEvent()
}

type CoreEvent struct{}

func (CoreEvent) embedCoreEvent() {}
