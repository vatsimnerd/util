package pubsub

import (
	"sync"
)

type InitialNotifier func(sub Subscription)

type Provider struct {
	dataReady bool
	subsLock  sync.RWMutex
	subs      map[string]Subscription
	notifier  InitialNotifier
}

func NewProvider() *Provider {
	return &Provider{
		subs: make(map[string]Subscription),
	}
}

func (p *Provider) SetInitialNotifier(notifier InitialNotifier) {
	p.notifier = notifier
}

func (p *Provider) SetDataReady(val bool) {
	p.dataReady = val
}

func (p *Provider) Subscribe(chSize int) Subscription {
	p.subsLock.Lock()
	defer p.subsLock.Unlock()
	sub := makeSubscription(chSize)
	p.subs[sub.id] = sub

	if p.dataReady && p.notifier != nil {
		p.notifier(sub)
	}
	return sub
}

func (p *Provider) Unsubscribe(sub Subscription) {
	p.subsLock.Lock()
	defer p.subsLock.Unlock()
	delete(p.subs, sub.id)
	close(sub.ch)
}

func (p *Provider) Notify(update Update) {
	p.subsLock.RLock()
	defer p.subsLock.RUnlock()
	for _, sub := range p.subs {
		sub.Send(update)
	}
}

func (p *Provider) Fin() {
	p.Notify(updateFin)
}
