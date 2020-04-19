package memo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
)

const (
	defaultDirName  = ".memo"
	defaultFileName = "memo.data"
)

type memo struct {
	store  map[string][]byte
	logger logger.Logger
}

func New(logger logger.Logger) (db.Client, error) {
	client := &memo{
		store:  map[string][]byte{},
		logger: logger,
	}

	err := client.LoadData()
	if err != nil {
		return nil, err
	}
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

func (this *memo) Close() error {
	return this.SaveData()
}

func (this *memo) LoadData() error {
	filePath := filepath.Join(os.Getenv("HOME"), defaultDirName, defaultFileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &this.store)
	if err != nil {
		return fmt.Errorf("error unmarshaling data from %s: %v", filePath, err)
	}

	return nil
}

func (this *memo) SaveData() error {
	memoDir := filepath.Join(os.Getenv("HOME"), defaultDirName)
	err := os.Mkdir(memoDir, 0700)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error mkdir %s: %v", memoDir, err)
	}

	filePath := filepath.Join(memoDir, defaultFileName)

	data, err := json.Marshal(this.store)
	if err != nil {
		return fmt.Errorf("error marshaling memo data: %v", err)
	}

	err = ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}

	return nil
}
