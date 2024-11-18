package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Coupon represents a coupon document in the database.
type Coupon struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Discount    float64            `json:"discount" bson:"discount"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

// CreateCouponParams is used to parse request data for creating a coupon.
type CreateCouponParams struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Discount    float64 `json:"discount" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

// UpdateCouponParams is used to parse request data for updating a coupon.
type UpdateCouponParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Discount    float64 `json:"discount"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

// CouponResponse is the structure returned in the API response for a coupon.
type CouponResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Discount    float64 `json:"discount"`
	CreatedAt   string  `json:"created_at"`
	ModifiedAt  string  `json:"modified_at"`
}
