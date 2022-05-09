package model

import "time"

type PointRequest struct {
	Cluster   string    `json:"cluster"`
	Type      string    `json:"type"`
	StartTime time.Time `json:"timerange_start"`
	EndTime   time.Time `json:"timerange_end"`
	// timequant
}
