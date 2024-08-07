package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type errorStruct struct {
	Error   string `json:"error,omitempty"`
	Details string `json:"errorDetails,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, dataStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	marshal, err := json.Marshal(dataStruct)
	if err != nil {
		newError := errorStruct{Error: http.StatusText(http.StatusInternalServerError)}
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
	errStruct := errorStruct{Error: http.StatusText(statusCode), Details: errorDetails}
	WriteJSON(w, statusCode, errStruct)
}
