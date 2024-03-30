package entity

import (
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID           int64        `db:"id" json:"id"`
	Email        string       `db:"email" json:"email"`
	PasswordHash string       `db:"password_hash" json:"-"`
	FirstName    string       `db:"first_name" json:"first_name"`
	LastName     string       `db:"last_name" json:"last_name"`
	Role         UserRoleEnum `db:"role" json:"role"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time   `db:"deleted_at" json:"-"`
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

type Product struct {
	ID           int64            `db:"id" json:"id"`
	Title        string           `db:"title" json:"title"`
	Subtitle     string           `db:"subtitle" json:"subtitle"`
	Description  string           `db:"description" json:"description"`
	Handle       string           `db:"handle" json:"handle"`
	IsPublished  bool             `db:"is_published" json:"is_published"`
	CollectionID *int64           `db:"collection_id" json:"collection_id"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time       `db:"deleted_at" json:"deleted_at"`
	Metadata     *json.RawMessage `db:"metadata" json:"metadata"`
}

type ProductFull struct {
}

type ProductCollection struct {
	ID        int64      `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Handle    string     `db:"handle" json:"handle"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ProductCategoryProduct struct {
	ProductID  int64 `db:"product_id"`
	CategoryID int64 `db:"category_id"`
}

type ProductImage struct {
	ProductID int64 `db:"product_id"`
	ImageID   int64 `db:"image_id"`
}

type ProductPrice struct {
	ID               int64      `db:"id"`
	ProductVariantID int64      `db:"product_variant_id"`
	CurrencyID       string     `db:"currency_id"`
	Amount           int        `db:"amount"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
}

type Price struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currency_code"`
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

type Discount struct {
	ID         int64            `db:"id" json:"id"`
	Code       string           `db:"code" json:"code"`
	IsActive   bool             `db:"is_active" json:"is_active"`
	Type       DiscountTypeEnum `db:"type" json:"type"`
	UsageLimit int              `db:"usage_limit" json:"usage_limit"`
	UsageCount int              `db:"usage_count" json:"usage_count"`
	StartsAt   time.Time        `db:"starts_at" json:"starts_at"`
	EndsAt     *time.Time       `db:"ends_at" json:"ends_at"`
	CreatedAt  time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time       `db:"deleted_at" json:"deleted_at"`
	Value      int              `db:"value" json:"value"`
}

type PaymentProvider struct {
	ID      string `db:"id"`
	Enabled bool   `db:"enabled"`
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
	DiscountID       int64             `db:"discount_id"`
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
	ID        int64     `db:"id"`
	ProductID int64     `db:"product_id"`
	Title     string    `db:"title"`
	Inventory int       `db:"inventory"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NotificationProvider struct {
	ID      string `db:"id"`
	Enabled bool   `db:"enabled"`
}

type Notification struct {
	ID           int64     `db:"id"`
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

type Invite struct {
	ID        int64        `db:"id"`
	Email     string       `db:"email"`
	Role      UserRoleEnum `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	ExpiresAt time.Time    `db:"expires_at"`
	Accepted  bool         `db:"accepted"`
	DeletedAt *time.Time   `db:"deleted_at"`
	Token     string       `db:"token"`
}

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDatabase        = errors.New("database error")
)
