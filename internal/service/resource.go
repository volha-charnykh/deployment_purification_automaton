package service

import (
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/db"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/model"
)

type Resource interface {
	Add(resources []model.Resource) ([]model.Resource, error)
	ListPoints(points model.PointRequest) ([]model.Resource, error)
}

type resource struct {
	db db.Resource
}

func NewResource(dao db.Resource) Resource {
	return resource{dao}
}

func (r resource) Add(res []model.Resource) ([]model.Resource, error) {
	rr := make([]model.Resource, len(res))
	for i, resource := range res {
		resp, err := r.db.Add(resource)
		if err != nil {
			return nil, err
		}
		rr[i] = resp
	}
	return rr, nil
}

func (r resource) ListPoints(req model.PointRequest) ([]model.Resource, error) {
	return r.db.ListForClusterFromInterval(req.Cluster, req.Type, req.StartTime, req.EndTime)
}
