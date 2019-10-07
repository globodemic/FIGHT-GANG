package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

//ResponsePlayer struct
type ResponsePlayer struct {
	ID             primitive.ObjectID `json:"id"  bson:"_id,omitempty"`
	Alias          string             `json:"alias" bson:"alias,omitempty"`
	Experience     int                `json:"exp" bson:"exp,omitempty"`
	Level          int                `json:"level" bson:"level,omitempty"`
	ExperienceNext int                `json:"expNext" bson:"expNext,omitempty"`
	Hits           int                `json:"hits" bson:"hits,omitempty"`
	Stamina        int                `json:"stamina" bson:"stamina,omitempty"`
	MaxStamina     int                `json:"maxStamina" bson:"maxStamina,omitempty"`
}
