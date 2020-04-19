package memo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
)

const (
	defaultDirName            = ".memo"
	defaultFileName           = "memo.data"
	defaultTimestampFormat    = time.RFC3339
	defaultSaveToDiskInterval = time.Hour
)

type memo struct {
	store   map[string][]byte
	logger  logger.Logger
	closeCh chan struct{}
}

func New(logger logger.Logger) (db.Client, error) {
	client := &memo{
		store:   map[string][]byte{},
		logger:  logger,
		closeCh: make(chan struct{}),
	}

	err := client.loadData()
	if err != nil {
		return nil, err
	}

	go client.periodicSave(defaultSaveToDiskInterval)

	return client, nil
}

func (this *memo) Get(key string) ([]byte, bool) {
	data, exist := this.store[key]
	if !exist {
		return nil, false
	}
	return data, true
}

func (this *memo) Set(key string, data []byte) bool {
	_, exist := this.store[key]
	if !exist {
		this.store[key] = data
		return false
	}
	this.store[key] = data
	return true
}

func (this *memo) GetKeys() []string {
	keys := make([]string, 0, len(this.store))
	for key := range this.store {
		keys = append(keys, key)
	}
	return keys
}

func (this *memo) Close() error {
	close(this.closeCh)
	return this.saveData()
}

func (this *memo) loadData() error {
	filePath := filepath.Join(os.Getenv("HOME"), defaultDirName, defaultFileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			this.logger.Printf("existing memo data not found at %s, skipping load from disk", filePath)
			return nil
		}
		return fmt.Errorf("error reading memo data from disk at %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &this.store)
	if err != nil {
		return fmt.Errorf("error unmarshaling memo data from %s: %v", filePath, err)
	}

	this.logger.Println("successfully loaded memo data from disk")
	return nil
}

func (this *memo) saveData() error {
	memoDir := filepath.Join(os.Getenv("HOME"), defaultDirName)
	err := os.Mkdir(memoDir, 0700)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error mkdir %s from saving memo data: %v", memoDir, err)
	}

	this.store["create_timestamp"] = []byte(time.Now().Format(defaultTimestampFormat))
	data, err := json.Marshal(this.store)
	if err != nil {
		return fmt.Errorf("error marshaling from saving memo data: %v", err)
	}

	filePath := filepath.Join(memoDir, defaultFileName)
	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing file %s from saving memo data: %v", filePath, err)
	}

	this.logger.Println("successfully saved memo data to disk")
	return nil
}

func (this *memo) periodicSave(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-this.closeCh:
			return
		case <-ticker.C:
			err := this.saveData()
			if err != nil {
				this.logger.Errorf("error from periodic save of memo data to disk: %v", err)
			}
		}
	}
}
