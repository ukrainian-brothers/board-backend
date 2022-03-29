package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type errorStruct struct {
	Error int `json:"error,omitempty"`
	Details string `json:"errorDetails,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, dataStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(dataStruct)
	w.WriteHeader(statusCode)
	if err != nil {
		newError := errorStruct{Error: http.StatusInternalServerError}
		bytes, _ := json.Marshal(newError)

		_, _ = w.Write(bytes)
		log.WithError(err).Error("marshaling http output failed")
	}
	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).Error("writing http output failed")
	}
}

func WriteError(w http.ResponseWriter, statusCode int, errorDetails string) {
	errStruct := errorStruct{Error: statusCode, Details: errorDetails}
	WriteJSON(w, statusCode, errStruct)
}
