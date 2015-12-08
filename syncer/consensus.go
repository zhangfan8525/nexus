package syncer

import (
	"path"
	"path/filepath"
	"pudding/datastructs"
	"pudding/myutils"
	"strings"
	"sync"
	"time"
)

var (
	TheConsensus *Consensus
	votesNeeded  int
)

type Consensus struct {
	sync.RWMutex
	nexus      *Nexus
	candidates chan datastructs.SyncFeedback
	m          map[StorageUnitId]int
	validated  map[StorageUnitId]bool
}

func NewConsensus(n *Nexus) *Consensus {
	return &Consensus{
		nexus:      n,
		candidates: make(chan datastructs.SyncFeedback, 1000),
		m:          make(map[StorageUnitId]int),
		validated:  make(map[StorageUnitId]bool),
	}
}

func (consensus *Consensus) MajorityPassAndVeto() {
	for sfb := range consensus.candidates {
		suid := GenStorageUnitId(sfb.Channel, sfb.Name, sfb.Revision)
		consensus.Lock()
		if _, ok := consensus.validated[suid]; !ok {
			if sfb.Valid {
				consensus.m[suid]++
				if consensus.m[suid] >= votesNeeded {
					consensus.validated[suid] = true
					c := consensus.nexus.LookupChannel(sfb.Channel)
					if c == nil {
						consensus.Unlock()
						return
					}
					su := c.LookupStorageUnit(sfb.Name, sfb.Revision)
					if su == nil {
						consensus.Unlock()
						return
					}
					su.Valid = true
					consensus.nexus.toBeSaved <- su
					myutils.MyLogger.Printf("Consensus passed: %+v", sfb)
				}
			} else {
				consensus.validated[suid] = false
				c := consensus.nexus.LookupChannel(sfb.Channel)
				if c == nil {
					consensus.Unlock()
					return
				}
				myutils.MyLogger.Printf("Consensus failed: %+v", sfb)
				c.RemoveStorageUnit(sfb.Name, sfb.Revision)
				ext := filepath.Ext(sfb.Name)
				woExt := strings.TrimSuffix(sfb.Name, ext)
				loc := path.Join(consensus.nexus.SyncRoot, sfb.Channel, woExt+"_r"+zeroPaddingLeft(3, sfb.Revision)+ext)
				consensus.nexus.toBeRemoved <- loc
				// TODO: load failed, need to get previous valid su
				// touch last su in the list?
				su := c.LookupLatestStorageUnit(sfb.Name)
				if su == nil {
					consensus.Unlock()
					return
				}
				su.TimeStamp = time.Now().Unix()
				myutils.MyLogger.Printf("Timestamp updated: %+v", su)
			}
		}
		consensus.Unlock()
	}
}
