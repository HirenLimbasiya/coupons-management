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
