package model

import "github.com/luispinto23/movieexample/metadata/pkg/model"

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata,omitempty"`
}
