package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/n4wei/memo/lib/logger"
	"github.com/n4wei/memo/model"
)

const (
	defaultDirName            = ".memo"
	defaultFileName           = "memo.data"
	defaultSaveToDiskInterval = time.Hour
)

type UserData struct {
	User  *model.User            `json:"user"`
	Memos map[string]*model.Memo `json:"memos"`
}

type Store struct {
	Users map[string]*UserData `json:"users"`
}

type cache struct {
	store   *Store
	logger  logger.Logger
	closeCh chan struct{}
}

func New(logger logger.Logger) (*cache, error) {
	cache := &cache{
		store:   &Store{Users: map[string]*UserData{}},
		logger:  logger,
		closeCh: make(chan struct{}),
	}

	err := cache.readFromDisk()
	if err != nil {
		return nil, err
	}

	go cache.periodicWriteToDisk(defaultSaveToDiskInterval)

	return cache, nil
}

func (this *cache) Close() error {
	close(this.closeCh)
	return this.writeToDisk()
}

func (this *cache) readFromDisk() error {
	filePath := filepath.Join(os.Getenv("HOME"), defaultDirName, defaultFileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			this.logger.Printf("no file at '%s', skipping read from disk", filePath)
			return nil
		}
		return fmt.Errorf("read from disk error, read from '%s': %v", filePath, err)
	}

	err = json.Unmarshal(data, &this.store)
	if err != nil {
		return fmt.Errorf("read from disk error, unmarshal data from '%s': %v", filePath, err)
	}

	this.logger.Println("successfully read from disk")
	return nil
}

func (this *cache) writeToDisk() error {
	dataDir := filepath.Join(os.Getenv("HOME"), defaultDirName)
	err := os.Mkdir(dataDir, 0700)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("write to disk error, mkdir '%s': %v", dataDir, err)
	}

	data, err := json.Marshal(this.store)
	if err != nil {
		return fmt.Errorf("write to disk error, marshal data: %v", err)
	}

	filePath := filepath.Join(dataDir, defaultFileName)
	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("write to disk error, write to '%s': %v", filePath, err)
	}

	this.logger.Println("successfully wrote to disk")
	return nil
}

func (this *cache) periodicWriteToDisk(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-this.closeCh:
			return
		case <-ticker.C:
			err := this.writeToDisk()
			if err != nil {
				this.logger.Errorf("periodic write to disk error: %v", err)
			}
		}
	}
}
