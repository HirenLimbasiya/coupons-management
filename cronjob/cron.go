package cronjob

import (
	"context"
	"coupons-management/db"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type CouponUpdater struct {
	Store *db.Store
}

func (c *CouponUpdater) UpdateExpiredCoupons() {
	ctx := context.Background()

	coupons, err := c.Store.Coupon.GetActiveCoupons(ctx)
	if err != nil {
		log.Printf("Error fetching active coupons: %v", err)
		return
	}

	now := time.Now()
	for _, coupon := range coupons {
		if coupon.ExpiresAt.Before(now) {
			err := c.Store.Coupon.UpdateCouponStatus(ctx, coupon.ID.Hex(), "Expired")
			if err != nil {
				log.Printf("Failed to update coupon %s: %v", coupon.ID.Hex(), err)
			} else {
				log.Printf("Coupon %s marked as expired", coupon.ID.Hex())
			}
		}
	}
}

func StartCouponCron(updater *CouponUpdater) {
	go func() {
		c := cron.New()
		_, err := c.AddFunc("@hourly", func() {
			log.Println("Running hourly job to update expired coupons...")
			updater.UpdateExpiredCoupons()
		})
		if err != nil {
			log.Fatalf("Error starting cron job: %v", err)
		}

		c.Start()
		log.Println("Coupon cron job started in a separate goroutine...")

		select {}
	}()
}
