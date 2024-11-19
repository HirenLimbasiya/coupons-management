package db

import (
	"context"
	"coupons-management/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	Coupon CouponStore
}

type CouponStore interface {
	GetCouponByID(ctx context.Context, couponID string) (*types.Coupon, error)
	GetAllCoupons(ctx context.Context) ([]types.Coupon, error)
	CreateCoupon(ctx context.Context, coupon types.CreateCouponParams) (*types.Coupon, error)
	UpdateCoupon(ctx context.Context, couponID string, coupon types.UpdateCouponParams) error
	DeleteCoupon(ctx context.Context, couponID string) error
	GetCouponsByType(ctx context.Context, couponID string) ([]types.Coupon, error)
}

func NewCouponStore(client *mongo.Client, dbName string) CouponStore {
	return &MongoCouponStore{
		client: client,
		dbName: dbName,
		collection: client.Database(dbName).Collection("coupons"),
	}
}

type MongoCouponStore struct {
	client *mongo.Client
	dbName string
	collection *mongo.Collection
}


func (s *MongoCouponStore) GetCouponByID(ctx context.Context, couponID string) (*types.Coupon, error) {

	objectID, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return nil, err 
	}

	filter := bson.M{
		"_id": objectID,
	}

	var coupon types.Coupon
	err = s.collection.FindOne(ctx, filter).Decode(&coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (s *MongoCouponStore) GetAllCoupons(ctx context.Context) ([]types.Coupon, error) {
	cursor, err := s.collection.Find(ctx, bson.M{}) // Empty filter to retrieve all coupons
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var coupons []types.Coupon
	for cursor.Next(ctx) {
		var coupon types.Coupon
		if err := cursor.Decode(&coupon); err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}
func (s *MongoCouponStore) GetCouponsByType(ctx context.Context, couponType string) ([]types.Coupon, error) {
	collection := s.client.Database(s.dbName).Collection("coupons")
	filter := bson.M{"type": couponType}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var coupons []types.Coupon
	for cursor.Next(ctx) {
		var coupon types.Coupon
		if err := cursor.Decode(&coupon); err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (s *MongoCouponStore) CreateCoupon(ctx context.Context, coupon types.CreateCouponParams) (*types.Coupon, error) {

	result, err := s.collection.InsertOne(ctx, coupon)
	if err != nil {
		return nil, err
	}

	return &types.Coupon{
		ID:          result.InsertedID.(primitive.ObjectID),
		Type: 		 coupon.Type,
		Details: coupon.Details,
        Description: coupon.Description,
        CreatedAt:   coupon.CreatedAt,
        ModifiedAt:  coupon.ModifiedAt,
		Status:      coupon.Status,
		ExpiresAt:   coupon.ExpiresAt,
	}, nil
}

func (s *MongoCouponStore) UpdateCoupon(ctx context.Context, couponID string, coupon types.UpdateCouponParams) error {
		objectID, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}
	update := map[string]interface{}{
		"$set": coupon,
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoCouponStore) DeleteCoupon(ctx context.Context, couponID string) error {
	collection := s.client.Database(s.dbName).Collection("coupons")
			objectID, err := primitive.ObjectIDFromHex(couponID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objectID,
	}

	_, err = collection.DeleteOne(ctx, filter)
	return err
}
