package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/n4wei/memo/db"
	"github.com/n4wei/memo/lib/logger"
	"github.com/n4wei/memo/model"
)

type userHandler struct {
	dbClient db.Client
	logger   logger.Logger
}

func NewUserHandler(dbClient db.Client, logger logger.Logger) *userHandler {
	return &userHandler{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (this *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		this.handleGet(w, r)
	case http.MethodPut:
		this.handlePut(w, r)
	default:
		handleBadRequest(w, r, errors.New("invalid request method"))
	}
}

func (this *userHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	users := this.dbClient.GetAllUsers()

	data, err := json.Marshal(users)
	if err != nil {
		handleServerError(w, r, fmt.Errorf("error marshaling users data: %v", err))
		return
	}

	addStandardResponseHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (this *userHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := handleRequestBody(r, user)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	user.CreateTimestamp = fmt.Sprintf("%v", time.Now().UTC().Unix())
	if ok := this.dbClient.AddUser(user.UserId, user); !ok {
		handleBadRequest(w, r, fmt.Errorf("user_id '%s' already exists, please choose a different user_id", user.UserId))
		return
	}

	addStandardResponseHeaders(w)
	w.WriteHeader(http.StatusCreated)
}
