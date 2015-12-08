package syncer

import (
	"os"
	"path"
	"path/filepath"
	"pudding/myutils"
	"strconv"
	"strings"
)

const (
	OkResponse    = `{"status": 0}`
	ErrorResponse = `{"status": 1}`
)

func zeroPaddingLeft(l, n int) string {
	s := strconv.Itoa(n)
	if len(s) >= l {
		return s
	} else {
		z := l - len(s)
		return strings.Repeat("0", z) + s
	}
}

func storageNameAndRevision(fn string) (sn string, rev int) {
	ext := filepath.Ext(fn)
	woExt := strings.TrimSuffix(fn, ext)
	l := len(woExt)
	// example_r002.xml
	if l < 5 {
		sn, rev = fn, 1
	} else if woExt[l-5:l-3] == "_r" {
		sn = woExt[:l-5] + ext
		rev, _ = strconv.Atoi(string(woExt[l-3:]))
	} else {
		sn, rev = fn, 1
	}
	return
}

func iterateSyncRoot(p string, fi os.FileInfo, e error) error {
	if p == TheNexus.SyncRoot {
		return nil
	}
	if !fi.IsDir() {
		return nil
	}
	err := filepath.Walk(p, loadChannelFolder)
	if err == filepath.SkipDir {
		return nil
	} else {
		return err
	}
}

func loadChannelFolder(p string, fi os.FileInfo, e error) error {
	if fi.IsDir() {
		return nil
	}
	_, err := TheNexus.LoadStorageUnitFromFile(p, path.Base(path.Dir(p)))
	if err != nil {
		myutils.MyLogger.Println(err.Error())
	}
	return err
}
