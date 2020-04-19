package list

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
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
	default:
		handleBadRequest(w, fmt.Errorf("invalid request method"))
	}
}

func (this *handler) handleGet(w http.ResponseWriter, r *http.Request) {
	keys := this.dbClient.GetKeys()
	data, err := json.Marshal(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"error marshaling keys": "%v"}`, err)))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func handleBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
}
