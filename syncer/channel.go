package syncer

import (
	"fmt"
	"pudding/myutils"
	"sort"
	"sync"
)

type Channel struct {
	Name  string
	Store map[string][]*StorageUnit
	lock  sync.RWMutex
}

func (c *Channel) String() (str string) {
	str = fmt.Sprintf("Channel name: %s | ", c.Name)
	for fileName, storageUnit := range c.Store {
		str += fmt.Sprintf("{Storage name: %s, Storage units: %s}", fileName, storageUnit)
	}
	return
}

func (c *Channel) SaveStorages(rootDir string) (err error) {
	for _, stores := range c.Store {
		for _, su := range stores {
			if su.Valid {
				err = su.SaveToFile(rootDir)
				if err != nil {
					myutils.MyLogger.Println(err.Error())
					return
				} else {
					myutils.MyLogger.Printf("Save storage unit: %+v\n", su)
				}
			}
		}
	}
	return nil
}

func (c *Channel) AddStorageUnit(s *StorageUnit) {
	c.lock.Lock()
	defer c.lock.Unlock()
	myutils.MyLogger.Printf("Add storage unit: %+v\n", s)
	c.Store[s.Name] = append(c.Store[s.Name], s)
}

func (c *Channel) RemoveStorageUnit(n string, r int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, exists := c.Store[n]; !exists {
		return
	} else {
		for i, su := range s {
			if su.Revision == r {
				myutils.MyLogger.Printf("Remove storage unit: %+v\n", su)
				s, s[len(s)-1] = append(s[:i], s[i+1:]...), nil
			}
		}
	}
}

func (c *Channel) RemoveByPredicate(n string, pred StorageUnitPredicate) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, exists := c.Store[n]; !exists {
		return
	} else {
		var passIndices []int
		for i, su := range s {
			if su != nil {
				if pred(su) {
					passIndices = append(passIndices, i)
				} else {
					myutils.MyLogger.Printf("Remove storage unit: %+v\n", su)
				}
			}
		}
		for i := 0; i < len(passIndices); i++ {
			s[i] = s[passIndices[i]]
		}
		for i := len(passIndices); i < len(s); i++ {
			s[i] = nil
		}
	}
}

func (c *Channel) RemoveByPredicateAll(pred StorageUnitPredicate) {
	for n, _ := range c.Store {
		c.RemoveByPredicate(n, pred)
	}
}

func (c *Channel) LookupStorageUnit(n string, r int) *StorageUnit {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if storageUnits, exists := c.Store[n]; !exists {
		return nil
	} else {
		for _, su := range storageUnits {
			if su == nil {
				continue
			}
			if su.Revision == r {
				return su
			}
		}
	}
	return nil
}

func (c *Channel) LookupLatestValidStorageUnit(n string) *StorageUnit {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if storageUnits, exists := c.Store[n]; !exists {
		return nil
	} else {
		if len(storageUnits) == 0 {
			return nil
		}
		for i := len(storageUnits) - 1; i >= 0; i-- {
			if storageUnits[i] == nil {
				continue
			}
			if storageUnits[i].Valid {
				return storageUnits[i]
			}
		}
	}
	return nil
}

func (c *Channel) LookupLatestStorageUnit(n string) *StorageUnit {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if storageUnits, exists := c.Store[n]; !exists {
		return nil
	} else {
		if len(storageUnits) == 0 {
			return nil
		}
		for i := len(storageUnits) - 1; i >= 0; i-- {
			if storageUnits[i] == nil {
				continue
			}
			return storageUnits[i]
		}
	}
	return nil
}

func (c *Channel) LookupUpdates(ts int64) (res []*StorageUnit) {
	if ts == 0 {
		for n, _ := range c.Store {
			su := c.LookupLatestStorageUnit(n)
			if su != nil {
				res = append(res, su)
			}
		}
	} else {
		for n, _ := range c.Store {
			su := c.LookupLatestStorageUnit(n)
			if su != nil && su.TimeStamp > ts-1 {
				res = append(res, su)
			}
		}
	}
	return
}

func (c *Channel) SortByRevision() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, storageUnits := range c.Store {
		sort.Sort(ByRevision(storageUnits))
	}
}

func (c *Channel) SortByTimeStamp() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, storageUnits := range c.Store {
		sort.Sort(ByTimeStamp(storageUnits))
	}
}
