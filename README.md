# Coupon Management System
This project is a coupon management system built using Go, the Fiber framework, and MongoDB. It supports multiple coupon types, including:

- Cart-wise: Discounts applied to the total cart value based on thresholds.
- Product-wise: Discounts applied to specific products in the cart.
- BxGy (Buy X Get Y): Buy specified quantities of products and get others for free, with repetition limits.
## Installation
1. Clone the repository:
```bash
git clone https://github.com/HirenLimbasiya/job-scheduler.git
```
2. Install dependencies:
```bash
go mod tidy
```
3. Set up the .env file: Create a .env file in the root directory and add the following configuration::
```bash
MONGO_URI=mongodb://localhost:27017
DB_NAME=coupon_system
PORT=8080
```
## Running the Application
1. Start the server::
```bash
go run main.go
```
2. The server will run on http://localhost:8080.


#### API Documentation
Document all API endpoints with examples.


## API Endpoints for Coupons

### 1. Create a Coupon
- Method: `POST`
- Endpoint: `/coupons`
- Payload:
```json
{
    "type": "cart-wise",
    "details": {
      "discount": 10,
      "threshold": 100
    }
}
```
- Response: 
```json
  {
    "coupon": {
        "id": "673bf444f42e23bae0eec36f",
        "type": "cart-wise",
        "details": {
            "discount": 10,
            "threshold": 100,
            "product_id": 0,
            "repetition_limit": 0,
            "buy_products": null,
            "get_products": null
        },
        "description": "",
        "created_at": "2024-11-19T07:43:24.254419+05:30",
        "modified_at": "2024-11-19T07:43:24.254419+05:30"
    }
  }
```

- **Other Coupons-Type Payload**: 

for type product-wise

```json
  {
    "type": "product-wise",
    "details": {
      "product_id": 1,
      "discount": 20
    }
  }
```
for type bxgy
```json
  {
    "type": "bxgy",
    "details": {
      "buy_products": [
        {
          "product_id": 1,
          "quantity": 3
        },
        {
          "product_id": 2,
          "quantity": 3
        }
      ],
      "get_products": [
        {
          "product_id": 3,
          "quantity": 1
        }
      ],
      "repition_limit": 2
    }
  }
```

### 2. Get Coupon By Id
- Method: `GET`
- Endpoint: `/coupons/:id`
- Response: 
```json
{
    "coupon": {
        "id": "673bf444f42e23bae0eec36f",
        "type": "cart-wise",
        "details": {
            "discount": 10,
            "threshold": 100,
            "product_id": 0,
            "repetition_limit": 0,
            "buy_products": null,
            "get_products": null
        },
        "description": "",
        "created_at": "2024-11-19T02:13:24.254Z",
        "modified_at": "2024-11-19T02:13:24.254Z"
    }
}
```

### 3. Get All Coupons
- Method: `GET`
- Endpoint: `/coupons`
- Response: 
```json
{
    "coupons": [{
        "id": "673bf444f42e23bae0eec36f",
        "type": "cart-wise",
        "details": {
            "discount": 10,
            "threshold": 100,
            "product_id": 0,
            "repetition_limit": 0,
            "buy_products": null,
            "get_products": null
        },
        "description": "",
        "created_at": "2024-11-19T02:13:24.254Z",
        "modified_at": "2024-11-19T02:13:24.254Z"
    }]
}
```

### 4. Update Coupon By ID
- Method: `PUT`
- Endpoint: `/coupons/:id`
- Payload: 
```json
  {
    "details": (According to type of update)
  }
```
- Response: 
```json
 {
    "message": "Coupon updated successfully"
 }
```

### 5. Delete Coupon By ID
- Method: `DELETE`
- Endpoint: `/coupons/:id`
- Response: 
```json
  {
    "message": "Coupon deleted successfully"
  }
```
## API Endpoints for Cart Operations

### 1. Get Applicable Coupons List For Cart
- Method: `POST`
- Endpoint: `/applicable-coupons`
- Payload:
```json
{
  "cart": {
    "items": [
      {
        "product_id": 1,
        "quantity": 6,
        "price": 50
      },
      {
        "product_id": 2,
        "quantity": 3,
        "price": 30
      },
      {
        "product_id": 3,
        "quantity": 2,
        "price": 25
      }
    ]
  }
}
```

- Response: 
```json
{
    "applicable_coupons": [
        {
            "coupon_id": "673bf7caf42e23bae0eec370",
            "type": "cart-wise",
            "discount": 44
        },
        {
            "coupon_id": "673b7c01bf5a40c7da13d0ab",
            "type": "product-wise",
            "discount": 60
        },
        {
            "coupon_id": "673b828f94ef11794030dd57",
            "type": "bxgy",
            "discount": 50
        }
    ]
}
```

### 2. Apply Coupon on Cart
- Method: `POST`
- Endpoint: `/apply-coupon/:id`
- Payload:
```json
{
  "cart": {
    "items": [
      {
        "product_id": 1,
        "quantity": 6,
        "price": 50
      },
      {
        "product_id": 2,
        "quantity": 3,
        "price": 30
      },
      {
        "product_id": 3,
        "quantity": 2,
        "price": 25
      }
    ]
  }
}
```
- Response: 
```json
{
    "updated_cart": {
        "final_price": 390,
        "items": [
            {
                "product_id": 1,
                "quantity": 6,
                "price": 50,
                "total_discount": 0
            },
            {
                "product_id": 2,
                "quantity": 3,
                "price": 30,
                "total_discount": 0
            },
            {
                "product_id": 3,
                "quantity": 4,
                "price": 25,
                "total_discount": 50
            }
        ],
        "total_discount": 50,
        "total_price": 440
    }
}
```

## Limitations
- Coupon Validation:
only one coupon allow at a time, stacking not allow
- Only on cart-wise type allow if multiple than it consider first one
- Dynamic coupon expiry not implemented yet 
- The BXGY coupons have a repetition_limit which is respected in applying free products but may not handle complex scenarios well

## Edge Cases
- Negative Price after Discount
- Product Quantities Not Matching the Coupon Requirement

## Possible Future Implementations
- Support for Stacking Discounts
- Coupon Expiry and Time Validation
- Product wise coupon apply only certain quantity
- Coupon Usage History
- Support for Percentage Off Coupons
- Coupon Prioritization
- Discount on Shipping

## Author

This project is created and maintained by **Hiren Limbasiya**.
You can explore more of my work on my [Portfolio](https://www.hirenlimbasiya.com/).