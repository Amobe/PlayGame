package domain

type Aggregator interface {
	ID() string
	Events() []Event
	Version() int
	embedCoreAggregator()
}

type CoreAggregator struct {
	version int
	events  []Event
}

func (a *CoreAggregator) Apply(new bool, events ...Event) {
	defer func() {
		if !new {
			a.version += len(events)
		}
	}()
	if new {
		a.events = append(a.events, events...)
	}
}

func (a CoreAggregator) Events() []Event {
	res := a.events
	a.events = nil
	return res
}

func (a CoreAggregator) Version() int {
	return a.version
}

func (a CoreAggregator) embedCoreAggregator() {}

type Event interface {
	Name() string
	embedCoreEvent()
}

type CoreEvent struct{}

func (CoreEvent) embedCoreEvent() {}
