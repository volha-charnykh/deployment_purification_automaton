package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/messages"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/mapper"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/model"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/service"
)

const AddResourceURL = "/cluster_agent_push_state"

type Resource struct {
	service service.Resource
}

func NewResource(svc service.Resource) Resource {
	return Resource{svc}
}

func (h Resource) HandleAddRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(fmt.Errorf(messages.ErrUnsupportedHttpMethod, r.Method), w)
		return
	}

	req, err := h.parseAddRequest(r.Body)
	if err != nil {
		writeErrorResponse(fmt.Errorf("Error while parsing request: %v", err), w)
		return
	}

	resp, err := h.service.Add(req)
	if err != nil {
		writeErrorResponse(fmt.Errorf("Error while handling request: %v", err), w)
		return
	}
	rr, err := prepareResponse(resp)
	if err != nil {
		writeErrorResponse(fmt.Errorf("Error while preparing response: %v", err), w)
		return
	}
	writeResponse(http.StatusOK, rr, w)

}

func (h Resource) parseAddRequest(body io.ReadCloser) ([]model.Resource, error) {
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf(messages.ErrCloseBody, err)
		}
	}()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf(messages.ErrReqRead, err)
	}
	data := &model.ResourceRequest{}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, fmt.Errorf(messages.ErrReqParse, err)
	}
	req, err := mapper.ParseResourceRequest(data)
	if err != nil {
		return nil, fmt.Errorf(messages.ErrReqParse, err)
	}
	return req, nil
}
