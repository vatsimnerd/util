package pubsub

import "testing"

func TestPubSub(t *testing.T) {
	const oType ObjectType = 1
	p := NewProvider()
	s := p.Subscribe(10)

	select {
	case upd := <-s.Updates():
		t.Errorf("unexpected update on channel: %v", upd)
	default:
	}

	p.Notify(Update{UType: UpdateTypeSet, OType: oType, Obj: "test"})
	select {
	case upd := <-s.Updates():
		if upd.UType != UpdateTypeSet {
			t.Errorf("unexpected update type %v, expected %v", upd.UType, UpdateTypeSet)
		}

		if upd.OType != oType {
			t.Errorf("unexpected object type %v, expected %v", upd.OType, oType)
		}

		obj, ok := upd.Obj.(string)
		if !ok {
			t.Error("expected string type")
		}

		if obj != "test" {
			t.Errorf("expected string 'test' in Obj, got %s", obj)
		}
	default:
		t.Errorf("expected update on channel")
	}
}

func TestNoInitialNotifier(t *testing.T) {
	const oType ObjectType = 1
	p := NewProvider()
	p.Notify(Update{UType: UpdateTypeSet, OType: oType, Obj: "test"})

	s := p.Subscribe(10)

	select {
	case upd := <-s.Updates():
		t.Errorf("unexpected update on channel: %v", upd)
	default:
	}
}

func TestInitialNotifier(t *testing.T) {
	const oType ObjectType = 1
	p := NewProvider()
	p.SetInitialNotifier(func(sub Subscription) {
		sub.Send(Update{UType: UpdateTypeSet, OType: oType, Obj: "test"})
	})
	p.SetDataReady(true)

	s := p.Subscribe(10)

	select {
	case upd := <-s.Updates():
		if upd.UType != UpdateTypeSet {
			t.Errorf("unexpected update type %v, expected %v", upd.UType, UpdateTypeSet)
		}

		if upd.OType != oType {
			t.Errorf("unexpected object type %v, expected %v", upd.OType, oType)
		}

		obj, ok := upd.Obj.(string)
		if !ok {
			t.Error("expected string type")
		}

		if obj != "test" {
			t.Errorf("expected string 'test' in Obj, got %s", obj)
		}
	default:
		t.Errorf("expected update on channel")
	}
}
