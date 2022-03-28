package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type errorStruct struct {
	Error HttpError `json:"error"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, dataStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(dataStruct)
	w.WriteHeader(statusCode)
	if err != nil {
		newError := errorStruct{Error: InternalError}
		bytes, _ := json.Marshal(newError)

		_, _ = w.Write(bytes)
		log.WithError(err).Error("marshaling http output failed")
	}
	_, err = w.Write(marshal)
	if err != nil {
		log.WithError(err).Error("writing http output failed")
	}
}

func WriteError(w http.ResponseWriter, error HttpError) {
	errStruct := errorStruct{Error: error}
	WriteJSON(w, errorsMap[error], errStruct)
}

type HttpError string

const (
	InternalError       HttpError = "internal_error"
	Unauthorized                  = "unauthorised"
	BadRequest                    = "bad_request"
	UnprocessableEntity           = "unprocessable_entity"
	InvalidPayload                = "unprocessable_entity"
)

var errorsMap = map[HttpError]int{
	InternalError:       500,
	Unauthorized:        401,
	BadRequest:          400,
	UnprocessableEntity: 422,
}