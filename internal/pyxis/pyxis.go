package pyxis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"
)

const (
	apiVersion = "v1"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type pyxisClient struct {
	APIToken  string
	ProjectID string
	Client    HTTPClient
	PyxisHost string
}

func (p *pyxisClient) getPyxisURL(path string) string {
	return fmt.Sprintf("https://%s/%s/%s", p.PyxisHost, apiVersion, path)
}

func (p *pyxisClient) getPyxisGraphqlURL() string {
	return fmt.Sprintf("https://%s/graphql/", p.PyxisHost)
}

func NewPyxisClient(pyxisHost string, apiToken string, projectID string, httpClient HTTPClient) *pyxisClient {
	return &pyxisClient{
		APIToken:  apiToken,
		ProjectID: projectID,
		Client:    httpClient,
		PyxisHost: pyxisHost,
	}
}

type CertImage struct {
	ID                     string           `json:"_id,omitempty"`
	Certified              bool             `json:"certified"`
	Deleted                bool             `json:"deleted" default:"false"`
	DockerImageDigest      string           `json:"docker_image_digest,omitempty"`
	DockerImageID          string           `json:"docker_image_id,omitempty"`
	ImageID                string           `json:"image_id,omitempty"`
	ISVPID                 string           `json:"isv_pid,omitempty"` // required
	ParsedData             *ParsedData      `json:"parsed_data,omitempty"`
	Architecture           string           `json:"architecture" default:"amd64"`
	RawConfig              string           `json:"raw_config,omitempty"`
	Repositories           []Repository     `json:"repositories,omitempty"`
	SumLayerSizeBytes      int64            `json:"sum_layer_size_bytes,omitempty"`
	UncompressedTopLayerID string           `json:"uncompressed_top_layer_id,omitempty"`
	FreshnessGrades        []FreshnessGrade `json:"freshness_grades,omitempty"`
}

type ParsedData struct {
	Architecture           string   `json:"architecture,omitempty"`
	Command                string   `json:"command,omitempty"`
	Comment                string   `json:"comment,omitempty"`
	Container              string   `json:"container,omitempty"`
	Created                string   `json:"created,omitempty"`
	DockerVersion          string   `json:"docker_version,omitempty"`
	ImageID                string   `json:"image_id,omitempty"`
	Labels                 []Label  `json:"labels,omitempty"` // required
	Layers                 []string `json:"layers,omitempty"` // required
	OS                     string   `json:"os,omitempty"`
	Ports                  string   `json:"ports,omitempty"`
	Size                   int64    `json:"size,omitempty"`
	UncompressedLayerSizes []Layer  `json:"uncompressed_layer_sizes,omitempty"`
}

type Repository struct {
	Published          bool   `json:"published" default:"false"`
	PushDate           string `json:"push_date,omitempty"` // time.Now
	Registry           string `json:"registry,omitempty"`
	Repository         string `json:"repository,omitempty"`
	Tags               []Tag  `json:"tags,omitempty"`
	ManifestListDigest string `json:"manifest_list_digest,omitempty"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type FreshnessGrade struct {
	Grade     string
	StartDate time.Time
	EndDate   time.Time
}

type Tag struct {
	AddedDate string `json:"added_date,omitempty"` // time.Now
	Name      string `json:"name,omitempty"`
}

type Layer struct {
	LayerID string `json:"layer_id"`
	Size    int64  `json:"size_bytes"`
}

// FindImagesByDigest uses an unauthenticated call to find_images() graphql function, and will
// return a slice of CertImages. It accepts a slice of image digests. The query return is then
// packed into the slice of CertImages.
func (p *pyxisClient) FindImagesByDigest(ctx context.Context, digests []string) ([]CertImage, error) {
	if len(digests) == 0 {
		return nil, fmt.Errorf("no digests specified")
	}
	// our graphQL query
	var query struct {
		FindImages struct {
			// Additional fields for return should be added here
			ContainerImage []struct {
				ID                graphql.String  `graphql:"_id"`
				Certified         graphql.Boolean `graphql:"certified"`
				DockerImageDigest graphql.String  `graphql:"docker_image_digest"`
			} `graphql:"data"`
			Error struct {
				Status graphql.Int    `graphql:"status"`
				Detail graphql.String `graphql:"detail"`
			} `graphql:"error"`
			Total graphql.Int
			Page  graphql.Int
			// filter to make sure we get exact results
		} `graphql:"find_images(filter: {docker_image_digest:{in:$digests}})"`
	}

	graphqlDigests := make([]graphql.String, len(digests))
	for idx, digest := range digests {
		graphqlDigests[idx] = graphql.String(digest)
	}
	// variables to feed to our graphql filter
	variables := map[string]interface{}{
		"digests": graphqlDigests,
	}

	// make our query
	httpClient, ok := p.Client.(*http.Client)
	if !ok {
		return nil, fmt.Errorf("client could not be used as http.Client")
	}
	client := graphql.NewClient(p.getPyxisGraphqlURL(), httpClient)

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("error while executing find_images query: %v", err)
	}

	images := make([]CertImage, len(query.FindImages.ContainerImage))
	for idx, image := range query.FindImages.ContainerImage {
		images[idx] = CertImage{
			ID:                string(image.ID),
			Certified:         bool(image.Certified),
			DockerImageDigest: string(image.DockerImageDigest),
		}
	}

	return images, nil
}
