package pubsub

type (
	UpdateType int
	ObjectType int

	Update struct {
		UType UpdateType
		OType ObjectType
		Obj   interface{}
	}
)

const (
	UpdateTypeSet UpdateType = iota
	UpdateTypeDelete
	UpdateTypeFin

	ObjectTypeNone ObjectType = 0
)

var (
	updateFin = Update{UType: UpdateTypeFin, OType: ObjectTypeNone, Obj: nil}
)

func MakeUpdates[T any](set map[string]T, del map[string]T, oType ObjectType) []Update {
	updates := make([]Update, 0)
	for _, v := range set {
		update := Update{
			UType: UpdateTypeSet,
			OType: oType,
			Obj:   v,
		}
		updates = append(updates, update)
	}
	for _, v := range del {
		update := Update{
			UType: UpdateTypeDelete,
			OType: oType,
			Obj:   v,
		}
		updates = append(updates, update)
	}
	return updates
}
