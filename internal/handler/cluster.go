package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/messages"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/model"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/service"
)

const ClusterGetPointsURL = "/cluster_get_points"

type Cluster struct {
	service service.Resource
}

func NewCluster(svc service.Resource) Cluster {
	return Cluster{svc}
}

func (h Cluster) HandleGetPointsRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(fmt.Errorf(messages.ErrUnsupportedHttpMethod, r.Method), w)
		return
	}

	req, err := h.parseGetPointsRequest(r.Body)
	if err != nil {
		writeErrorResponse(fmt.Errorf("Error while parsing request: %v", err), w)
		return
	}

	resp, err := h.service.ListPoints(req)
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

func (h Cluster) parseGetPointsRequest(body io.ReadCloser) (model.PointRequest, error) {
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf(messages.ErrCloseBody, err)
		}
	}()
	data := model.PointRequest{}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return data, fmt.Errorf(messages.ErrReqRead, err)
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return data, fmt.Errorf(messages.ErrReqParse, err)
	}
	return data, nil
}
