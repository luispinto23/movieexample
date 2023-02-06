package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/luispinto23/movieexample/movie/internal/gateway"
	"github.com/luispinto23/movieexample/pkg/discovery"
	"github.com/luispinto23/movieexample/rating/pkg/model"
)

// Gateway defines an HTTP gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetServiceURL(ctx context.Context) (string, error) {
	// fetch the available addresses
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return "", err
	}

	// chose a random address from those available
	url := "http://" + addrs[rand.Intn(len(addrs))] + "rating"
	return url, nil
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	url, err := g.GetServiceURL(ctx)
	if err != nil {
		return 0, err
	}
	log.Printf("Calling rating service. Request GET: " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", string(recordType))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	// fetch the available addresses
	url, err := g.GetServiceURL(ctx)
	if err != nil {
		return err
	}

	log.Printf("calling rating service. request get: " + url)
	req, err := http.NewRequest(http.MethodPut, url, nil)

	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
