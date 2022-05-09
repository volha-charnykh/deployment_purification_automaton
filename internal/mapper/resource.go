package mapper

import (
	"fmt"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/messages"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/model"
)

func ParseResourceRequest(rr *model.ResourceRequest) ([]model.Resource, error) {
	if rr == nil || len(rr.Resources) == 0 {
		return nil, fmt.Errorf("empty resource")
	}
	resources := make([]model.Resource, len(rr.Resources))
	for i, r := range rr.Resources {
		kind, name, namespace, err := parseResourceInternal(r)
		if err != nil {
			return nil, err
		}
		resources[i] = model.Resource{
			Resource:      r,
			Cluster:       rr.Cluster,
			Timestamp:     rr.Time,
			ResourceType:  kind,
			ResourceName:  name,
			NamespaceName: namespace,
			AppName:       rr.Application,
			AppVersion:    rr.Version,
			GitHash:       rr.GitHash,
			State:         rr.State,
		}
	}
	return resources, nil
}

func parseResourceInternal(r map[string]interface{}) (string, string, string, error) {
	k, ok := r["kind"]
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	kind, ok := k.(string)
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	meta, ok := r["metadata"]
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	m, ok := meta.(map[string]interface{})
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	nn, ok := m["namespace"]
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	namespace, ok := nn.(string)
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	n, ok := m["name"]
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}
	name, ok := n.(string)
	if !ok {
		return "", "", "", fmt.Errorf(messages.ErrInvalidResource)
	}

	return kind, name, namespace, nil
}
