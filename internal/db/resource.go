package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/messages"
	"github.com/Infoblox-CTO/deployment_purification_automaton/internal/model"
)

const (
	ClusterTypeK8       = "k8"
	ClusterTypeManifest = "manifest"
)

type Resource interface {
	ListForClusterFromInterval(cluster, resourceType string, time1, time2 time.Time) ([]model.Resource, error)
	Add(resource model.Resource) (model.Resource, error)
}

type resource struct {
	db *sql.DB
}

func NewResource(db *sql.DB) Resource {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "resources"
	(
		resource        JSONB NOT NULL,
		cluster      TEXT NOT NULL,
		cluster_type      TEXT NOT NULL,
		git_hash     TEXT,
		ts TIMESTAMPTZ NOT NULL,
		resource_type   TEXT  NOT NULL,
		resource_name   TEXT  NOT NULL,
		namespace_name   TEXT  NOT NULL,
		app_name   TEXT,
		app_version   TEXT,
		state   TEXT,
		UNIQUE (cluster,cluster_type,resource_type,resource_name,namespace_name,ts)
	);
	CREATE INDEX resources_git_hash_idx ON resources (git_hash) WHERE cluster_type='` + ClusterTypeManifest + `';`)
	if err != nil {
		log.Fatalf(messages.ErrDbCreateTable, "resources", err)
	}

	return resource{db}
}

func (r resource) ListForClusterFromInterval(cluster, resourceType string, time1, time2 time.Time) ([]model.Resource, error) {
	var resources []model.Resource

	rows, err := r.db.Query("SELECT resource, cluster, cluster_type, git_hash, ts, resource_type, resource_name, namespace_name, app_name, app_version, state FROM resources WHERE cluster = $1 AND resource_type=$2 AND ts >= $3 AND ts <= $4", cluster, resourceType, time1, time2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var res model.Resource
		if err := rows.Scan(&res.Resource, &res.Cluster, &res.ClusterType, &res.GitHash, &res.Timestamp, &res.ResourceType, &res.ResourceName, &res.NamespaceName, &res.AppName, &res.AppVersion, &res.State); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(messages.ErrDbRead, err)
	}
	return resources, nil
}

func (r resource) Add(resource model.Resource) (model.Resource, error) {
	var res model.Resource
	row := r.db.QueryRow("INSERT INTO resources (resource, cluster, cluster_type, git_hash, ts, resource_type, resource_name, namespace_name, app_name, app_version, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING resource, cluster, git_hash, ts, resource_type, resource_name, namespace_name, app_name, app_version, state",
		resource.Resource, resource.Cluster, resource.ClusterType, resource.GitHash, resource.Timestamp, resource.ResourceType, resource.ResourceName, resource.NamespaceName, resource.AppName, resource.AppVersion, resource.State)
	if err := row.Scan(&res.Resource, &res.Cluster, &res.ClusterType, &res.GitHash, &res.Timestamp, &res.ResourceType, &res.ResourceName, &res.NamespaceName, &res.AppName, &res.AppVersion, &res.State); err != nil {
		return res, fmt.Errorf(messages.ErrDbAdd, err)
	}
	return res, nil
}
