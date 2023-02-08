package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types.
type RecordType string

// Existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// RatingEventType defines the type of a rating event.
type RatingEventType string

// Rating event types.
const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	RecordID   string      `json:"record_id,omitempty"`
	RecordType string      `json:"record_type,omitempty"`
	UserID     UserID      `json:"user_id,omitempty"`
	Value      RatingValue `json:"value,omitempty"`
}

// RatingEvent defines an event containing rating information.
type RatingEvent struct {
	UserID     UserID
	RecordID   RecordID
	RecordType RecordType
	Value      RatingValue
	EventType  RatingEventType
}
