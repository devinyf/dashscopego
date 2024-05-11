package qwen

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"os"
	"time"
)

const MaxFileCacheLifeTime = time.Hour*48 - time.Minute*5

var gFileCacheMgr *FileCacheMgr

type FileCache struct {
	URL        string
	UploadTime int64
}

type FileCacheMgr struct {
	MapFiles map[string]*FileCache
}

func (mgr *FileCacheMgr) hash(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func (mgr *FileCacheMgr) Cache(buf []byte, url string) {
	key := mgr.hash(buf)

	mgr.MapFiles[key] = &FileCache{
		URL:        url,
		UploadTime: time.Now().Unix(),
	}
}

func (mgr *FileCacheMgr) Get(buf []byte) string {
	key := mgr.hash(buf)

	cache, isok := mgr.MapFiles[key]
	if isok {
		curtime := time.Now().Unix()
		if curtime-cache.UploadTime <= int64(MaxFileCacheLifeTime) {
			return cache.URL
		}

		delete(mgr.MapFiles, key)
	}

	return ""
}

func (mgr *FileCacheMgr) Save(fn string) error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	enc.Encode(mgr)

	return os.WriteFile(fn, buf.Bytes(), 0644)
}

func LoadFileCacheMgr(fn string) error {
	buf, err := os.ReadFile(fn)
	if err != nil {
		gFileCacheMgr = &FileCacheMgr{
			MapFiles: make(map[string]*FileCache),
		}

		return nil
	}

	gFileCacheMgr = &FileCacheMgr{}

	dec := gob.NewDecoder(bytes.NewReader(buf))

	err = dec.Decode(gFileCacheMgr)
	if err != nil {
		return err
	}

	return nil
}

func SaveFileCacheMgr(fn string) error {
	if gFileCacheMgr != nil {
		return gFileCacheMgr.Save(fn)
	}

	return nil
}
