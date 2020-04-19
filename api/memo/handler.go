package memo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
	"github.com/n4wei/memo/model"
)

type handler struct {
	dbClient db.Client
	logger   logger.Logger
}

func NewHandler(dbClient db.Client, logger logger.Logger) http.Handler {
	return &handler{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (this *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		this.handleGet(w, r)
	case http.MethodPut:
		this.handlePut(w, r)
	default:
		handleBadRequest(w, fmt.Errorf("invalid request method"))
	}
}

func (this *handler) handleGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := pathParts[2]
	data, exist := this.dbClient.Get(key)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"key":"%s","value":%s}`, key, data)))
}

func (this *handler) handlePut(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w, fmt.Errorf("error reading request body: %v", err))
		return
	}

	var data model.Memo
	err = json.Unmarshal(body, &data)
	if err != nil {
		handleBadRequest(w, fmt.Errorf("error unmarshaling request body: %v", err))
		return
	}

	exist := this.dbClient.Set(data.Key, data.Value)
	if !exist {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
}
