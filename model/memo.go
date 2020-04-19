package model

import "encoding/json"

type Memo struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}
