package memory

import (
	"context"

	"github.com/luispinto23/movieexample/rating/internal/repository"
	"github.com/luispinto23/movieexample/rating/pkg/model"
)

// The following implementation is using a nested map to store all records inside it. If we didn’t define
// separate types, RatingID, RatingType, and UserID, it would be harder to understand the types
// of the keys in the map because we would be using primitives such as string and int, which are
// less self-descriptive.

// Repository defines a rating repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	ratings, ok := r.data[recordType][recordID]

	if !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return ratings, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)

	return nil
}
