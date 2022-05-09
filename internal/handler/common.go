package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func prepareResponse(data interface{}) ([]byte, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		return nil, err
	}
	return body.Bytes(), nil
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	log.Println("Error: ", err)

	writeResponse(http.StatusBadRequest, []byte(err.Error()), w)
	return
}

func writeResponse(statusCode int, body []byte, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println("Response parsing error: ", err)
		writeErrorResponse(err, w)
	}
	log.Println("Response body: ", string(body))
}
