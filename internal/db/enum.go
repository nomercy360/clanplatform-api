package db

type DiscountTypeEnum string

const (
	Percentage   DiscountTypeEnum = "percentage"
	Fixed        DiscountTypeEnum = "fixed"
	FreeShipping DiscountTypeEnum = "free_shipping"
)

type OrderStatusEnum string

const (
	Pending    OrderStatusEnum = "pending"
	Processing OrderStatusEnum = "processing"
	Shipped    OrderStatusEnum = "shipped"
	Delivered  OrderStatusEnum = "delivered"
	Cancelled  OrderStatusEnum = "cancelled"
)

type PaymentStatusEnum string

const (
	PaymentPending   PaymentStatusEnum = "pending"
	PaymentConfirmed PaymentStatusEnum = "confirmed"
	Refunded         PaymentStatusEnum = "refunded"
)
