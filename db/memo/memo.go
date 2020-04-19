package memo

import (
	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
)

type memo struct {
	store  map[string][]byte
	logger logger.Logger
}

func New(logger logger.Logger) db.Client {
	return &memo{
		store:  map[string][]byte{},
		logger: logger,
	}
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
