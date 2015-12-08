package syncer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"pudding/myutils"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultSyncRoot           = "sync_root"
	StorageUnitExpiryDuration = 2 * 60 * 60
)

var (
	ErrSameAsLast = errors.New("md5sum same as last file")
)

type StorageUnit struct {
	Name      string
	Channel   string
	Revision  int
	Valid     bool
	Md5sum    string // hex
	TimeStamp int64  // Unix time
	Data      []byte
}

type StorageUnitId string

func (s *StorageUnit) String() (str string) {
	str = fmt.Sprintf("{Name: %s, Channel: %s, Revision: %d, Valid: %t, Md5: %s, Modified time: %d, Data length: %d}", s.Name, s.Channel, s.Revision, s.Valid, s.Md5sum, s.TimeStamp, len(s.Data))
	return
}

func (s *StorageUnit) Id() StorageUnitId {
	return StorageUnitId(s.Channel + ":" + s.Name + ":" + strconv.Itoa(s.Revision))
}

func GenStorageUnitId(c, n string, r int) StorageUnitId {
	return StorageUnitId(c + ":" + n + ":" + strconv.Itoa(r))
}

func GetStorageUnitFields(id string) (channel, name string, revision int) {
	fields := strings.Split(id, ":")
	if len(fields) != 3 {
		return "", "", 0
	}
	channel = fields[0]
	name = fields[1]
	revision, err := strconv.Atoi(fields[2])
	if err != nil {
		revision = 1
		return
	}
	return
}

type StorageUnitPredicate func(*StorageUnit) bool

func StorageUnitBeforeExpiry(s *StorageUnit) bool {
	return time.Now().Unix()-s.TimeStamp < int64(StorageUnitExpiryDuration)
}

func StorageUnitValid(s *StorageUnit) bool {
	return s.Valid
}

func (s *StorageUnit) SaveToFile(rootDir string) (err error) {
	if !s.Valid {
		return fmt.Errorf("Trying to save invalid file %s", s.Name)
	}
	ext := filepath.Ext(s.Name)
	woExt := strings.TrimSuffix(s.Name, ext)
	loc := path.Join(rootDir, s.Channel, woExt+"_r"+zeroPaddingLeft(3, s.Revision)+ext)
	dir := filepath.Dir(loc)
	if ok, _ := myutils.Exists(dir); !ok {
		err := os.MkdirAll(dir, os.ModeDir|0755)
		if err != nil {
			myutils.MyLogger.Println("Counld not create " + dir + " directory: " + err.Error())
			return err
		}
	}
	return ioutil.WriteFile(loc, s.Data, 0666)
}

type ByRevision []*StorageUnit

func (by ByRevision) Len() int {
	return len(by)
}

func (by ByRevision) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}

func (by ByRevision) Less(i, j int) bool {
	return by[i].Revision < by[j].Revision
}

type ByTimeStamp []*StorageUnit

func (by ByTimeStamp) Len() int {
	return len(by)
}

func (by ByTimeStamp) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}

func (by ByTimeStamp) Less(i, j int) bool {
	return by[i].TimeStamp < by[j].TimeStamp
}

func init() {
	needFeedBack := false
	if needFeedBackStr, ok := myutils.Configs["sync_need_feedback"]; !ok {
		needFeedBack = false
	} else {
		if strings.ToLower(needFeedBackStr) == "true" {
			needFeedBack = true
		} else {
			needFeedBack = false
		}
	}
	if needFeedBack {
		if votesNeededStr, ok := myutils.Configs["consensus_votes"]; !ok {
			votesNeeded = 1
		} else {
			var err error
			votesNeeded, err = strconv.Atoi(votesNeededStr)
			if err != nil {
				votesNeeded = 1
				myutils.MyLogger.Println(err.Error())
			}
		}
	} else {
		votesNeeded = 0
	}

	TheNexus = NewNexus(Mediocre, DefaultSyncRoot, 200, needFeedBack)
	TheConsensus = NewConsensus(TheNexus)
	myutils.MyLogger.Println(fmt.Sprintf("Nexus is running in %s mode.", TheNexus.Mode))
	myutils.MyLogger.Println("Nexus need feedback?:", needFeedBack, "; votes needed:", votesNeeded)
	sr := TheNexus.SyncRoot
	if ok, _ := myutils.Exists(sr); !ok {
		err := os.Mkdir(sr, os.ModeDir|0755)
		if err != nil {
			myutils.MyLogger.Println("Counld not create " + sr + " directory: " + err.Error())
		}
	} else {
		filepath.Walk(sr, iterateSyncRoot)
	}
	// sort by revision so "latest" su will be the one with largest rev #
	for idx, c := range TheNexus.Channels {
		c.SortByRevision()
		myutils.MyLogger.Println("Loaded channel " + strconv.Itoa(idx) + ": " + fmt.Sprintf("%s", c))
	}

	go TheConsensus.MajorityPassAndVeto()

	go func() {
		ticker := time.Tick(1 * time.Hour)
		for {
			select {
			case su := <-TheNexus.toBeSaved:
				su.SaveToFile(TheNexus.SyncRoot)
			case loc := <-TheNexus.toBeRemoved:
				if ok, _ := myutils.Exists(loc); ok {
					//os.Remove(loc)
					// previously saved file becomes invalid
					os.Rename(loc, loc+".oblivion")
				}
			case <-ticker:
				// TODO: SUs that have not yet been ACKed will be purged, persist to disk?
				// should set a timeout for each and every su
				// can be done by comparing timestamp of su and NOW, purge(persist?) it if the span is larger than some duration
				if needFeedBack {
					TheNexus.PurgeExpired()
				}
				//TheNexus.PurgeInvalid()
			}
		}
	}()
	/*
		TheNexus.SaveChannels("testsave")
		su1, _ := TheNexus.NewStorageUnitFromData(&datastructs.SyncPut{
			Name:    "FromDataTest",
			Channel: "test",
			Content: "Test save from data",
		})
		su1.Valid = true
		su2, _ := TheNexus.NewStorageUnitFromData(&datastructs.SyncPut{
			Name:    "FromDataTest",
			Channel: "test",
			Content: "Test save from data, take 2",
		})
		su2.Valid = true
		TheNexus.toBeSaved <- &StorageUnit{
			Name:     "chantest",
			Channel:  "test",
			Revision: 12,
			Valid:    true,
			Data:     []byte("Hello World"),
		}
		TheNexus.toBeSaved <- &StorageUnit{
			Name:     "chantest",
			Channel:  "test",
			Revision: 13,
			Valid:    true,
			Data:     []byte("Hello World"),
		}
		select {
		case <-time.After(1 * time.Second):
			TheNexus.SaveChannels("testsave")
		}
	*/
}
