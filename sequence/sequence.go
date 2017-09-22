package sequence

import (
	"fmt"
	"sort"
	"sync"
)

var (
	sequencesMu sync.RWMutex
	sequences   = map[string]Sequence{}
)

type Sequence interface {
	NextSequence() (seq uint64, err error)
}

// GetSequence returns corresponding sequence instance with the specified
// sequenceType. If the specified sequenceType does not register itself, then
// err will be non nil and sequence will be nil. Else, err will be nil and
// sequence will be corresponding sequence instance.
func GetSequence(sequenceType string) (sequence Sequence, err error) {
	sequencesMu.RLock()
	defer sequencesMu.RUnlock()

	if value, ok := sequences[sequenceType]; ok {
		sequence = value
		return sequence, nil
	} else {
		return nil, fmt.Errorf("%v is not registered.", sequenceType)
	}
}

// Register makes a sequence generator available by the provided sequenceType.
// If Register is called twice with the same name or if driver is nil, it
// panics.
func Register(sequenceType string, sequence Sequence) {
	sequencesMu.Lock()
	defer sequencesMu.Unlock()

	if sequence == nil {
		panic("sequence: Registered sequence is nil")
	}

	if _, dup := sequences[sequenceType]; dup {
		panic("sequence: Register called twice for driver " + sequenceType)
	}

	sequences[sequenceType] = sequence
}

// Sequences returns a sorted list of the types of the registered sequences.
func Sequences() []string {
	sequencesMu.RLock()
	defer sequencesMu.RUnlock()

	var list []string
	for name := range sequences {
		list = append(list, name)
	}

	sort.Strings(list)
	return list
}
