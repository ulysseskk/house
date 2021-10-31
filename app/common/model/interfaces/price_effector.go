package interfaces

type PriceEffector interface {
	GetTotalPrice() uint64
	SetTotalPrice(price uint64)
	UnitPriceString() string
	GetUnitPrice() uint64
	SetUnitPrice(unitPrice uint64)
}
