package syncer

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"pudding/datastructs"
	"pudding/myutils"
	"sync"
	"time"
)

type NexusMode int

const (
	Mediocre       NexusMode = 0
	ShinyAndChrome NexusMode = 1
)

// TODO: make Nexus shiny and chrome
func (n NexusMode) String() (str string) {
	switch n {
	case Mediocre:
		str = "Mediocre"
	case ShinyAndChrome:
		str = "Shiny and Chrome"
	default:
		str = "Wild World"
	}
	return
}

type Nexus struct {
	Channels     []*Channel
	Mode         NexusMode
	SyncRoot     string
	maxChans     int
	lock         sync.RWMutex
	needFeedback bool
	toBeSaved    chan *StorageUnit
	toBeRemoved  chan string
}

func NewNexus(mode NexusMode, root string, max int, needFb bool) *Nexus {
	return &Nexus{
		Mode:         mode,
		SyncRoot:     root,
		maxChans:     max,
		needFeedback: needFb,
		toBeSaved:    make(chan *StorageUnit, 100),
		toBeRemoved:  make(chan string, 100),
	}
}

var (
	TheNexus           *Nexus
	ErrChannelOverflow = errors.New("Reached the maximum number of channels this nexus can handle")
)

func (nexus *Nexus) LoadStorageUnitFromFile(fullName, channelName string) (*StorageUnit, error) {
	if filepath.Ext(fullName) == ".oblivion" {
		// silently fade into oblivion
		return nil, nil
	}

	data, err := ioutil.ReadFile(fullName)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s, %s", fullName, err.Error())
	}
	sum := fmt.Sprintf("%x", md5.Sum(data))
	name, rev := storageNameAndRevision(filepath.Base(fullName))

	channel := nexus.LookupChannel(channelName)
	if channel == nil {
		channel = nexus.CreateNewChannel(channelName)
		if channel == nil {
			return nil, ErrChannelOverflow
		}
	}

	su := &StorageUnit{
		Name:      name,
		Channel:   channelName,
		Revision:  rev,
		Valid:     true, // only valid su can be stored
		Md5sum:    sum,
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}
	channel.AddStorageUnit(su)

	return su, nil
}

func (nexus *Nexus) NewStorageUnitFromData(sp *datastructs.SyncPut) (*StorageUnit, error) {
	data := []byte(sp.Content)
	sum := fmt.Sprintf("%x", md5.Sum(data))
	channel := nexus.LookupChannel(sp.Channel)
	var rev int
	if channel != nil {
		if s := channel.LookupLatestStorageUnit(sp.Name); s != nil {
			if sum != s.Md5sum {
				rev = s.Revision + 1
			} else {
				// touch last su in channel, so it will be retrived by the consumers again, might not be the desired behaviour
				s.TimeStamp = time.Now().Unix()
				myutils.MyLogger.Printf("Timestamp updated: %+v", s)
				return nil, ErrSameAsLast
			}
		} else {
			rev = 1
		}
	} else {
		channel = nexus.CreateNewChannel(sp.Channel)
		if channel == nil {
			return nil, ErrChannelOverflow
		}
		rev = 1
	}

	su := &StorageUnit{
		Name:      sp.Name,
		Channel:   sp.Channel,
		Revision:  rev,
		Valid:     false, // needs response from client
		Md5sum:    sum,
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}
	channel.AddStorageUnit(su)

	return su, nil
}

func (nexus *Nexus) SaveChannels(rootDir string) (err error) {
	nexus.lock.RLock()
	defer nexus.lock.RUnlock()
	for _, c := range nexus.Channels {
		err = c.SaveStorages(rootDir)
		if err != nil {
			myutils.MyLogger.Println(err.Error())
			return
		}
	}
	return nil
}

func (nexus *Nexus) PurgeInvalid() {
	myutils.MyLogger.Println("Trying to purge invalid storage units")
	nexus.lock.Lock()
	defer nexus.lock.Unlock()
	for _, c := range nexus.Channels {
		c.RemoveByPredicateAll(StorageUnitValid)
	}
}

func (nexus *Nexus) PurgeExpired() {
	myutils.MyLogger.Println("Trying to purge expired storage units")
	nexus.lock.Lock()
	defer nexus.lock.Unlock()
	for _, c := range nexus.Channels {
		c.RemoveByPredicateAll(StorageUnitBeforeExpiry)
	}
}

func (nexus *Nexus) LookupChannel(n string) *Channel {
	nexus.lock.RLock()
	defer nexus.lock.RUnlock()
	for _, c := range nexus.Channels {
		if c.Name == n {
			return c
		}
	}
	return nil
}

func (nexus *Nexus) CreateNewChannel(n string) *Channel {
	nexus.lock.Lock()
	defer nexus.lock.Unlock()
	if len(nexus.Channels) == nexus.maxChans {
		return nil
	}

	c := &Channel{
		Name:  n,
		Store: make(map[string][]*StorageUnit),
	}
	nexus.Channels = append(nexus.Channels, c)
	return c
}
