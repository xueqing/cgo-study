package object

import "sync"

// ID ...
type ID int32

var refs struct {
	sync.Mutex
	objs map[ID]interface{}
	next ID
}

func init() {
	refs.Lock()
	defer refs.Unlock()

	refs.objs = make(map[ID]interface{})
	refs.next = 1
}

// NewID ...
func NewID(obj interface{}) ID {
	refs.Lock()
	defer refs.Unlock()

	id := refs.next
	refs.next++

	refs.objs[id] = obj
	return id
}

// IsNil ...
func (id ID) IsNil() bool {
	return id == 0
}

// Get ...
func (id ID) Get() interface{} {
	refs.Lock()
	defer refs.Unlock()

	return refs.objs[id]
}

// Free ...
func (id *ID) Free() interface{} {
	refs.Lock()
	defer refs.Unlock()

	obj := refs.objs[*id]
	delete(refs.objs, *id)
	*id = 0

	return obj
}
