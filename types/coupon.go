package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Type 		string 			   `json:"type" bson:"type"`
	Details 	CouponDetails      `json:"details" bson:"details"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	Status      string             `json:"status" bson:"status"`
}

type CouponDetails struct {
	Discount  		float64 		  `json:"discount"`
	Threshold 		float64 		  `json:"threshold"`
	ProductID 		int 			  `json:"product_id" bson:"product_id"`
	RepetitionLimit int 			  `json:"repetition_limit" bson:"repetition_limit"`
	BuyProducts 	[]ProductQuantity `json:"buy_products" bson:"buy_products"`
	GetProducts 	[]ProductQuantity `json:"get_products" bson:"get_products"`
}
type CreateCouponParams struct {
	Type 		string 			   `json:"type" bson:"type"`
	Description string  		   `json:"description" validate:"required"`
	Details 	CouponDetails 	   `json:"details" bson:"details"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

type ProductQuantity struct {
    ProductID int `json:"product_id" bson:"product_id"`
    Quantity  int `json:"quantity"`
}


type UpdateCouponParams struct {
	Details 	CouponDetails 	   `json:"details" bson:"details"`
	Description string  		   `json:"description"`
	ModifiedAt  time.Time          `json:"modified_at" bson:"modified_at"`
}

type Cart struct {
	Cart CartData `json:"cart" bson:"cart"`
}

type CartData struct {
    Items []CartItem `json:"items" bson:"items"`
}
type CartItem struct {
	ProductID 	  int     `json:"product_id"`
	Quantity  	  int     `json:"quantity"`
	Price     	  float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount" bson:"total_discount"`
}

type ApplicableCoupon struct {
	CouponID string  `json:"coupon_id"`
	Type     string  `json:"type"`
	Discount float64 `json:"discount"`
}