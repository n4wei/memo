package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handleRequestBody(r *http.Request, dataModel interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}

	err = json.Unmarshal(body, dataModel)
	if err != nil {
		return fmt.Errorf("error unmarshaling request body: %v", err)
	}

	return nil
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	writeErrorResponse(w, r, err)
}

func handleNotFound(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusNotFound)
	if err != nil {
		writeErrorResponse(w, r, err)
	}
}

func handleServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	writeErrorResponse(w, r, err)
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"error": "%s %s %s\n%v"`, r.Proto, r.Method, r.URL.Path, err)))
}
