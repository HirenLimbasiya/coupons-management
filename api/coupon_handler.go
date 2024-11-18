package api

import (
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

	req.ModifiedAt = time.Now()
	err := h.store.Coupon.UpdateCoupon(c.Context(), couponID, req)
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

	// Fetch all cart-wise coupons
	coupons, err := h.store.Coupon.GetCouponsByType(c.Context(), "cart-wise")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching coupons",
		})
	}

	// Find applicable coupons
	var applicableCoupons []types.ApplicableCoupon
	if len(coupons) > 0 {
    coupon := coupons[0] // Get the first coupon

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

	productCoupons, err := h.store.Coupon.GetCouponsByType(c.Context(), "product-wise")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching product-wise coupons",
		})
	}

	// Find applicable product-wise coupons
	for _, coupon := range productCoupons {

		// Check if the coupon applies to any item in the cart
		for _, item := range req.Cart.Items {
			if item.ProductID == coupon.Details.ProductID {
				// Calculate the discount amount based on the quantity and discount
				discountAmount := (float64(item.Quantity) * item.Price * coupon.Details.Discount) / 100

				// Add the applicable product-wise coupon to the response
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

	// Find applicable bxgy coupons
	for _, coupon := range bxgyCoupons {
		// Extract details from the coupon
		buyProducts := coupon.Details.BuyProducts
		getProducts := coupon.Details.GuyProducts
		repetitionLimit := coupon.Details.RepetitionLimit

		// Calculate how many times the coupon can be applied based on the cart's product quantities
		var totalBuyQuantity int
		for _, buyProduct := range buyProducts {
			for _, item := range req.Cart.Items {
				if item.ProductID == buyProduct.ProductID {
					totalBuyQuantity += item.Quantity
				}
			}
		}

		// Calculate the possible repetitions of the coupon
		repetitions := totalBuyQuantity / buyProducts[0].Quantity // Use the first item in the buyProducts array
		if repetitions > repetitionLimit {
			repetitions = repetitionLimit
		}

		// Apply the coupon if it can be applied
		if repetitions > 0 {
			var discountAmount float64
			// For each repetition, calculate the discount for the free products
			for i := 0; i < repetitions; i++ {
				for _, getProduct := range getProducts {
					for _, item := range req.Cart.Items {
						if item.ProductID == getProduct.ProductID {
							// Add the price of the free product to the discount
							discountAmount += float64(getProduct.Quantity) * item.Price
						}
					}
				}
			}

			// Add the applicable BxGy coupon to the response
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

	// Fetch the coupon by ID
	coupon, err := h.store.Coupon.GetCouponByID(c.Context(), couponID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Coupon not found",
		})
	}

	// Initialize total discount and final cart price
	var totalDiscount float64
	var finalPrice float64
	for _, item := range req.Cart.Items {
		finalPrice += float64(item.Quantity) * item.Price
	}

	// Apply the coupon based on its type
	switch coupon.Type {
	case "cart-wise":
		// For cart-wise coupon, check if the threshold is met and apply the discount

		// If the total cart value exceeds the threshold, apply the discount
		if coupon.Details.Threshold <= finalPrice {
			discountAmount := (finalPrice * coupon.Details.Discount) / 100
			totalDiscount += discountAmount
			finalPrice -= discountAmount
		}

	case "product-wise":
		// Apply discount to the product if it exists in the cart
		for i, item := range req.Cart.Items {
			if item.ProductID == coupon.Details.ProductID {
				// Calculate the discount on the total price (quantity * price)
				discountAmount := (float64(item.Quantity) * item.Price * coupon.Details.Discount) / 100
				req.Cart.Items[i].TotalDiscount = discountAmount
				// Subtract the discount amount from the total price of the item
				req.Cart.Items[i].Price -= (discountAmount / float64(item.Quantity)) // Adjust price per item
				totalDiscount += discountAmount
				finalPrice -= discountAmount
			}
		}

	case "bxgy":
		// For BxGy (Buy X Get Y) coupon, apply based on the buy and get products
		buyProducts := coupon.Details.BuyProducts
		getProducts := coupon.Details.GuyProducts
		repetitionLimit := coupon.Details.RepetitionLimit

		// Calculate how many times the coupon can be applied based on the cart's product quantities
		var totalBuyQuantity int
		for _, buyProduct := range buyProducts {
			for _, item := range req.Cart.Items {
				if item.ProductID == buyProduct.ProductID {
					totalBuyQuantity += item.Quantity
				}
			}
		}

		// Calculate the possible repetitions of the coupon
		repetitions := totalBuyQuantity / buyProducts[0].Quantity // Use the first item in the buyProducts array
		if repetitions > repetitionLimit {
			repetitions = repetitionLimit
		}

		// Apply the coupon if it can be applied
		if repetitions > 0 {
			var discountAmount float64
			// Update the quantity and discount for get products (free products)
			for _, getProduct := range getProducts {
				for i, item := range req.Cart.Items {
					if item.ProductID == getProduct.ProductID {
						// Calculate how many free products we can add
						freeQuantity := getProduct.Quantity * repetitions
						req.Cart.Items[i].Quantity += freeQuantity
						// Add the price of the free product to the discount
						discountAmount += float64(freeQuantity) * item.Price
						req.Cart.Items[i].TotalDiscount = float64(freeQuantity) * item.Price
					}
				}
			}

			// Apply the discount and update the cart
			totalDiscount += discountAmount
			finalPrice -= discountAmount
		}
	}

	// Return the updated cart with total price and total discount applied
	return c.JSON(fiber.Map{
		"updated_cart": fiber.Map{
			"items": req.Cart.Items,
			"total_price": finalPrice + totalDiscount, // Add total discount to get the final price
			"total_discount": totalDiscount,
			"final_price": finalPrice,
		},
	})
}
