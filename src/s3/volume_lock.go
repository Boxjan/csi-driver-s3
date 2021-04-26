package s3

import (
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	volumeOperationAlreadyExistsFmt = "An operation with the given Volume ID %s already exists"
)

// VolumeLocks implements a map with atomic operations. It stores a set of all volume IDs
// with an ongoing operation.
type volumeLocks struct {
	locks sets.String
	mux   sync.Mutex
}

func newVolumeLocks() *volumeLocks {
	return &volumeLocks{
		locks: sets.NewString(),
	}
}

// TryAcquire tries to acquire the lock for operating on volumeID and returns true if successful.
// If another operation is already using volumeID, returns false.
func (vl *volumeLocks) TryAcquire(volumeID string) bool {
	vl.mux.Lock()
	defer vl.mux.Unlock()
	if vl.locks.Has(volumeID) {
		return false
	}
	vl.locks.Insert(volumeID)
	return true
}

func (vl *volumeLocks) Release(volumeID string) {
	vl.mux.Lock()
	defer vl.mux.Unlock()
	vl.locks.Delete(volumeID)
}
