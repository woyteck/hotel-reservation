package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	RoomID      primitive.ObjectID `bson:"roomID" json:"roomID"`
	NumPersons  int                `bson:"numPersons" json:"numPersons"`
	FromDate    time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate    time.Time          `bson:"tillDate" json:"tillDate"`
	IsCancelled bool               `bson:"isCancelled" json:"isCancelled"`
}
