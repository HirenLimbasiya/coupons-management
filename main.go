package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"coupons-management/api"
	"coupons-management/cronjob"
	"coupons-management/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoEndpoint := os.Getenv("MONGO_URL")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	store, err := NewStore(mongoEndpoint, dbName)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	app := fiber.New()

	couponHandler := api.NewCouponHandler(*store)

	// Define your routes
	app.Get("/coupons/:id", couponHandler.HandleGetCoupon)
	app.Get("/coupons", couponHandler.HandleGetAllCoupons)
	app.Post("/coupons", couponHandler.HandleCreateCoupon)
	app.Put("/coupons/:id", couponHandler.HandleUpdateCoupon)
	app.Delete("/coupons/:id", couponHandler.HandleDeleteCoupon)

	app.Post("/applicable-coupons", couponHandler.HandleGetApplicableCoupons)
	app.Post("/apply-coupon/:id", couponHandler.HandleApplyCoupon)

	couponUpdater := &cronjob.CouponUpdater{Store: store}
	cronjob.StartCouponCron(couponUpdater)

	app.Listen(":" + port)
}

func NewStore(mongoURI string, dbName string) (*db.Store, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	store := &db.Store{
		Coupon: db.NewCouponStore(client, dbName),
	}

	return store, nil
}

