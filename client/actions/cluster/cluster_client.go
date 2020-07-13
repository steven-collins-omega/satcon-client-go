package cluster

import (
	"errors"
	"net/http"

	"github.ibm.com/coligo/satcon-client/client"
)

// ClusterService is the interface used to perform all cluster-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ClusterService
type ClusterService interface {
	// RegisterCluster registers a new cluster under the specified organization ID.
	RegisterCluster(string, Registration, string) (*RegisterClusterResponseDataDetails, error)
	// ClustersByOrgID lists the clusters registered under the specified organization.
	ClustersByOrgID(string, string) (ClusterList, error)
	// DeleteClusterByClusterID deletes the specified cluster from the specified org,
	// including all resources under that cluster.
	DeleteClusterByClusterID(string, string, string) (*DeleteClustersResponseDataDetails, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	Endpoint   string
	HTTPClient client.HTTPClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient client.HTTPClient) (ClusterService, error) {
	if endpointURL == "" {
		return nil, errors.New("Must supply a valid endpoint URL")
	}

	if httpClient == nil {
		return &Client{
			Endpoint:   endpointURL,
			HTTPClient: http.DefaultClient,
		}, nil
	}

	return &Client{
		Endpoint:   endpointURL,
		HTTPClient: httpClient,
	}, nil
}