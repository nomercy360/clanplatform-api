package entity

import (
	"time"
)

type User struct {
	ID           int64        `db:"id"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	FirstName    string       `db:"first_name"`
	LastName     string       `db:"last_name"`
	Role         UserRoleEnum `db:"role"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    *time.Time   `db:"deleted_at"`
}

type Customer struct {
	ID         int64      `db:"id"`
	Email      string     `db:"email"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	HasAccount bool       `db:"has_account"`
}

type Image struct {
	ID        int64      `db:"id"`
	URL       string     `db:"url"`
	CreatedAt time.Time  `db:"created_at"`
	IsMain    bool       `db:"is_main"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Currency struct {
	Code   string `db:"code"`
	Name   string `db:"name"`
	Symbol string `db:"symbol"`
}

type Category struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ProductCategory struct {
	ProductID  int64 `db:"product_id"`
	CategoryID int64 `db:"category_id"`
}

type ProductImage struct {
	ProductID int64 `db:"product_id"`
	ImageID   int64 `db:"image_id"`
}

type ProductPrice struct {
	ID         int64      `db:"id"`
	ProductID  int64      `db:"product_id"`
	CurrencyID string     `db:"currency_id"`
	Price      int        `db:"price"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type Product struct {
	ID             int64      `db:"id"`
	Title          string     `db:"title"`
	Subtitle       string     `db:"subtitle"`
	Description    string     `db:"description"`
	Handle         string     `db:"handle"`
	Brand          string     `db:"brand"`
	Condition      string     `db:"condition"`
	Material       string     `db:"material"`
	Model          string     `db:"model"`
	CollectionID   int64      `db:"collection_id"`
	ProductionNote string     `db:"production_note"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type Address struct {
	ID         int64      `db:"id"`
	CustomerID int64      `db:"customer_id"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	Address    string     `db:"address"`
	City       string     `db:"city"`
	Province   string     `db:"province"`
	Country    string     `db:"country"`
	Zip        string     `db:"zip"`
	Phone      string     `db:"phone"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type Promotion struct {
	ID         int64            `db:"id"`
	Code       string           `db:"code"`
	IsActive   bool             `db:"is_active"`
	Type       DiscountTypeEnum `db:"type"`
	UsageLimit int              `db:"usage_limit"`
	EndsAt     time.Time        `db:"ends_at"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"`
	DeletedAt  *time.Time       `db:"deleted_at"`
}

type PaymentProvider struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Enabled   bool       `db:"enabled"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Payment struct {
	ID                int64             `db:"id"`
	OrderID           int64             `db:"order_id"`
	PaymentProviderID int64             `db:"payment_provider_id"`
	Status            PaymentStatusEnum `db:"status"`
	Amount            int               `db:"amount"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`
	DeletedAt         *time.Time        `db:"deleted_at"`
}

type Order struct {
	ID               int64             `db:"id"`
	CustomerID       int64             `db:"customer_id"`
	AddressID        int64             `db:"address_id"`
	PromotionID      int64             `db:"promotion_id"`
	Status           OrderStatusEnum   `db:"status"`
	PaymentStatus    PaymentStatusEnum `db:"payment_status"`
	TotalPrice       int               `db:"total_price"`
	CreatedAt        time.Time         `db:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at"`
	DeletedAt        *time.Time        `db:"deleted_at"`
	ShippingMethodID int64             `db:"shipping_method_id"`
	OrderNumber      string            `db:"order_number"`
}

type OrderItem struct {
	ID        int64      `db:"id"`
	OrderID   int64      `db:"order_id"`
	ProductID int64      `db:"product_id"`
	Quantity  int        `db:"quantity"`
	Price     int        `db:"price"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type ShippingMethod struct {
	ID           int64     `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Cost         float64   `db:"cost"`
	DeliveryTime string    `db:"delivery_time"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type ProductVariant struct {
	ID           int64              `db:"id"`
	ProductID    int64              `db:"product_id"`
	VariantName  ProductVariantEnum `db:"variant_name"`
	VariantValue string             `db:"variant_value"`
	SKU          string             `db:"sku"`
	Inventory    int                `db:"inventory"`
	CreatedAt    time.Time          `db:"created_at"`
	UpdatedAt    time.Time          `db:"updated_at"`
}

type Notification struct {
	ID           string    `db:"id"`
	EventName    string    `db:"event_name"`
	ResourceType string    `db:"resource_type"`
	ResourceID   string    `db:"resource_id"`
	CustomerID   string    `db:"customer_id"`
	To           string    `db:"to"`
	Data         string    `db:"data"`
	ProviderID   string    `db:"provider_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
