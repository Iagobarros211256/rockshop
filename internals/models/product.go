package models

//here  is modes to database, i wil eplce after with real db requests
//product gonna represent a generic product at sale
//obviously i wll add more aand extend it

type Prodduct struct {
	ID          uint   `json:"id"`
	SKU         string `json:"sku" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description, omitempty"`
	PriceCents  int64  `json:"price_cents" binding:"required"`
	Stock       int    `json:"stock"`
	Type        string `json:"type"`
}
