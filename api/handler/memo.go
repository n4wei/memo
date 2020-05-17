package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
	"github.com/n4wei/memo/model"
)

type memoHandler struct {
	dbClient db.Client
	logger   logger.Logger
}

func NewMemoHandler(dbClient db.Client, logger logger.Logger) *memoHandler {
	return &memoHandler{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (this *memoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		this.handleGet(w, r)
	case http.MethodPut:
		this.handlePut(w, r)
	default:
		handleBadRequest(w, r, fmt.Errorf("invalid request method"))
	}
}

func (this *memoHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	userId, err := parsePath(r.URL.Path)
	if err != nil {
		handleNotFound(w, r, nil)
		return
	}

	memos, exist := this.dbClient.GetAllUserMemos(userId)
	if !exist {
		handleNotFound(w, r, nil)
		return
	}

	data, err := json.Marshal(memos)
	if err != nil {
		handleServerError(w, r, fmt.Errorf("error marshaling memo data: %v", err))
		return
	}

	addStandardResponseHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (this *memoHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	userId, err := parsePath(r.URL.Path)
	if err != nil {
		handleNotFound(w, r, nil)
		return
	}

	memo := &model.Memo{}
	err = handleRequestBody(r, memo)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	memo.MemoId = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	memo.UserId = userId

	added := this.dbClient.AddUserMemo(userId, memo)
	if !added {
		handleNotFound(w, r, nil)
		return
	}

	addStandardResponseHeaders(w)
	w.WriteHeader(http.StatusCreated)
}

func parsePath(path string) (string, error) {
	path = path[1:]
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return "", errors.New("invalid path")
	}

	return parts[1], nil
}
