package api

import (
	"fmt"
	"net/http"
	"time"

	"coupons-management/db"
	"coupons-management/types"

	"github.com/gofiber/fiber/v2"
)

type CouponHandler struct {
	store db.Store
}

func NewCouponHandler(store db.Store) *CouponHandler {
	return &CouponHandler{store: store}
}

func (h *CouponHandler) HandleCreateCoupon(c *fiber.Ctx) error {
	var req types.CreateCouponParams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}
	if req.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ExpiresAt must be a future date",
		})
	}

	if err := validateCouponDetails(req.Details, req.Type); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.Status = "Active"
	req.CreatedAt = time.Now() 
	req.ModifiedAt = time.Now()

	coupon, err := h.store.Coupon.CreateCoupon(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating coupon"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"coupon": coupon})
}

func (h *CouponHandler) HandleGetCoupon(c *fiber.Ctx) error {
	couponID := c.Params("id")
	coupon, err := h.store.Coupon.GetCouponByID(c.Context(), couponID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
	}

	return c.JSON(fiber.Map{"coupon": coupon})
}

func (h *CouponHandler) HandleGetAllCoupons(c *fiber.Ctx) error {
	coupons, err := h.store.Coupon.GetAllCoupons(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving coupons"})
	}

	return c.JSON(fiber.Map{"coupons": coupons})
}

func (h *CouponHandler) HandleUpdateCoupon(c *fiber.Ctx) error {
	couponID := c.Params("id")

	var req types.UpdateCouponParams
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	coupon, err := h.store.Coupon.GetCouponByID(c.Context(), couponID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
	}

	if err := validateCouponDetails(req.Details, coupon.Type); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	req.ModifiedAt = time.Now()
	err = h.store.Coupon.UpdateCoupon(c.Context(), couponID, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating coupon"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Coupon updated successfully"})
}

func (h *CouponHandler) HandleDeleteCoupon(c *fiber.Ctx) error {
	couponID := c.Params("id")
	err := h.store.Coupon.DeleteCoupon(c.Context(), couponID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Coupon not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Coupon deleted successfully"})
}

func (h *CouponHandler) HandleGetApplicableCoupons(c *fiber.Ctx) error {
	var req types.Cart
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	totalCartValue := 0.0
	for _, item := range req.Cart.Items {
		totalCartValue += float64(item.Quantity) * item.Price
	}

	coupons, err := h.store.Coupon.GetCouponsByType(c.Context(), "cart-wise")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching coupons",
		})
	}

	var applicableCoupons []types.ApplicableCoupon
	if len(coupons) > 0 {
    coupon := coupons[0]
		if coupon.Status == "Active" && coupon.ExpiresAt.After(time.Now()) {
			threshold := coupon.Details.Threshold
			discount := coupon.Details.Discount
		
				if threshold <= totalCartValue {
					discountAmount := (totalCartValue * discount) / 100
					applicableCoupons = append(applicableCoupons, types.ApplicableCoupon{
						CouponID: coupon.ID.Hex(),
						Type:     coupon.Type,
						Discount: discountAmount,
					})
			
			}
		}
	}

	productCoupons, err := h.store.Coupon.GetCouponsByType(c.Context(), "product-wise")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching product-wise coupons",
		})
	}

	for _, coupon := range productCoupons {
		if coupon.Status == "Expired" || coupon.ExpiresAt.Before(time.Now()) {
			continue 
		}
		for _, item := range req.Cart.Items {
			if item.ProductID == coupon.Details.ProductID {
				discountAmount := (float64(item.Quantity) * item.Price * coupon.Details.Discount) / 100

				applicableCoupons = append(applicableCoupons, types.ApplicableCoupon{
					CouponID: coupon.ID.Hex(),
					Type:     coupon.Type,
					Discount: discountAmount,
				})
			}
		}
	}

	bxgyCoupons, err := h.store.Coupon.GetCouponsByType(c.Context(), "bxgy")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching bxgy coupons",
		})
	}

	for _, coupon := range bxgyCoupons {
		if coupon.Status == "Expired" || coupon.ExpiresAt.Before(time.Now()) {
			continue 
		}
		buyProducts := coupon.Details.BuyProducts
		getProducts := coupon.Details.GetProducts
		repetitionLimit := coupon.Details.RepetitionLimit

		var totalBuyQuantity int
		for _, buyProduct := range buyProducts {
			for _, item := range req.Cart.Items {
				if item.ProductID == buyProduct.ProductID {
					totalBuyQuantity += item.Quantity
				}
			}
		}

		repetitions := totalBuyQuantity / buyProducts[0].Quantity
		if repetitions > repetitionLimit {
			repetitions = repetitionLimit
		}

		if repetitions > 0 {
			var discountAmount float64
			for i := 0; i < repetitions; i++ {
				for _, getProduct := range getProducts {
					for _, item := range req.Cart.Items {
						if item.ProductID == getProduct.ProductID {
							discountAmount += float64(getProduct.Quantity) * item.Price
						}
					}
				}
			}

			applicableCoupons = append(applicableCoupons, types.ApplicableCoupon{
				CouponID: coupon.ID.Hex(),
				Type:     coupon.Type,
				Discount: discountAmount,
			})
		}
	}

	return c.JSON(fiber.Map{
		"applicable_coupons": applicableCoupons,
	})
}

func (h *CouponHandler) HandleApplyCoupon(c *fiber.Ctx) error {
	couponID := c.Params("id")

	var req types.Cart
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	coupon, err := h.store.Coupon.GetCouponByID(c.Context(), couponID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	if coupon.Status == "Expired" || coupon.ExpiresAt.Before(time.Now()) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Coupon has expired",
		})
	}

	var totalDiscount float64
	var finalPrice float64
	for _, item := range req.Cart.Items {
		finalPrice += float64(item.Quantity) * item.Price
	}

	switch coupon.Type {
	case "cart-wise":
		if coupon.Details.Threshold <= finalPrice {
			discountAmount := (finalPrice * coupon.Details.Discount) / 100
			totalDiscount += discountAmount
			finalPrice -= discountAmount
		}

	case "product-wise":
		for i, item := range req.Cart.Items {
			if item.ProductID == coupon.Details.ProductID {
				discountAmount := (float64(item.Quantity) * item.Price * coupon.Details.Discount) / 100
				req.Cart.Items[i].TotalDiscount = discountAmount
				req.Cart.Items[i].Price -= (discountAmount / float64(item.Quantity)) 
				totalDiscount += discountAmount
				finalPrice -= discountAmount
			}
		}

	case "bxgy":
		buyProducts := coupon.Details.BuyProducts
		getProducts := coupon.Details.GetProducts
		repetitionLimit := coupon.Details.RepetitionLimit

		var totalBuyQuantity int
		for _, buyProduct := range buyProducts {
			for _, item := range req.Cart.Items {
				if item.ProductID == buyProduct.ProductID {
					totalBuyQuantity += item.Quantity
				}
			}
		}

		repetitions := totalBuyQuantity / buyProducts[0].Quantity
		if repetitions > repetitionLimit {
			repetitions = repetitionLimit
		}

		if repetitions > 0 {
			var discountAmount float64
			for _, getProduct := range getProducts {
				for i, item := range req.Cart.Items {
					if item.ProductID == getProduct.ProductID {
						freeQuantity := getProduct.Quantity * repetitions
						req.Cart.Items[i].Quantity += freeQuantity
						discountAmount += float64(freeQuantity) * item.Price
						req.Cart.Items[i].TotalDiscount = float64(freeQuantity) * item.Price
					}
				}
			}

			totalDiscount += discountAmount
			finalPrice -= discountAmount
		}
	}

	return c.JSON(fiber.Map{
		"updated_cart": fiber.Map{
			"items": req.Cart.Items,
			"total_price": finalPrice + totalDiscount,
			"total_discount": totalDiscount,
			"final_price": finalPrice,
		},
	})
}

func validateCouponDetails(details types.CouponDetails, couponType string) error {
	switch couponType {
	case "cart-wise":
		if details.Threshold == 0 || details.Discount == 0 {
			return fmt.Errorf("for cart-wise type, both Threshold and Discount must be provided")
		}

	case "product-wise":
		if details.ProductID == 0 || details.Discount == 0 {
			return fmt.Errorf("for product-wise type, both ProductID and Discount must be provided")
		}

	case "bxgy":
		if len(details.BuyProducts) == 0 || len(details.GetProducts) == 0 || details.RepetitionLimit == 0 {
			return fmt.Errorf("for bxgy type, BuyProducts, GetProducts, and RepetitionLimit must be provided")
		}

	default:
		return fmt.Errorf("invalid coupon type. Only 'cart-wise', 'product-wise', and 'bxgy' are allowed")
	}
	return nil
}
