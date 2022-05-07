package pubsub

import "github.com/google/uuid"

type Subscription struct {
	id string
	ch chan Update
}

func makeSubscription(chSize int) Subscription {
	return Subscription{
		id: uuid.New().String(),
		ch: make(chan Update, chSize),
	}
}

func (s Subscription) Send(update Update) {
	s.ch <- update
}

func (s Subscription) Fin() {
	s.ch <- updateFin
}
