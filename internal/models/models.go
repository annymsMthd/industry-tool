package models

import "time"

type EveAsset struct {
	ItemID          int64  `json:"item_id"`
	IsBlueprintCopy bool   `json:"is_blueprint_copy"`
	IsSingleton     bool   `json:"is_singleton"`
	LocationFlag    string `json:"location_flag"`
	LocationID      int64  `json:"location_id"`
	LocationType    string `json:"location_type"`
	Quantity        int64  `json:"quantity"`
	TypeID          int64  `json:"type_id"`
}

type EveInventoryType struct {
	TypeID   int64
	TypeName string
	Volume   float64
	IconID   *int64
}

type Region struct {
	ID   int64
	Name string
}

type Constellation struct {
	ID       int64
	Name     string
	RegionID int64
}

type SolarSystem struct {
	ID              int64
	Name            string
	ConstellationID int64
	Security        float64
}

type Station struct {
	ID            int64
	Name          string
	SolarSystemID int64
	CorporationID int64
	IsNPC         bool
}

type Corporation struct {
	ID   int64
	Name string
}

type CorporationDivisions struct {
	Hanger map[int]string
	Wallet map[int]string
}

type StockpileMarker struct {
	UserID          int64   `json:"userId"`
	TypeID          int64   `json:"typeId"`
	OwnerType       string  `json:"ownerType"`
	OwnerID         int64   `json:"ownerId"`
	LocationID      int64   `json:"locationId"`
	ContainerID     *int64  `json:"containerId"`
	DivisionNumber  *int    `json:"divisionNumber"`
	DesiredQuantity int64   `json:"desiredQuantity"`
	Notes           *string `json:"notes"`
}

type MarketPrice struct {
	TypeID      int64
	RegionID    int64
	BuyPrice    *float64
	SellPrice   *float64
	DailyVolume *int64
	UpdatedAt   string
}

type Contact struct {
	ID              int64      `json:"id"`
	RequesterUserID int64      `json:"requesterUserId"`
	RecipientUserID int64      `json:"recipientUserId"`
	RequesterName   string     `json:"requesterName"`
	RecipientName   string     `json:"recipientName"`
	Status          string     `json:"status"`
	RequestedAt     time.Time  `json:"requestedAt"`
	RespondedAt     *time.Time `json:"respondedAt"`
}

type ContactPermission struct {
	ID              int64  `json:"id"`
	ContactID       int64  `json:"contactId"`
	GrantingUserID  int64  `json:"grantingUserId"`
	ReceivingUserID int64  `json:"receivingUserId"`
	ServiceType     string `json:"serviceType"`
	CanAccess       bool   `json:"canAccess"`
}

type ForSaleItem struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"userId"`
	TypeID            int64     `json:"typeId"`
	TypeName          string    `json:"typeName"`
	OwnerType         string    `json:"ownerType"`
	OwnerID           int64     `json:"ownerId"`
	OwnerName         string    `json:"ownerName"`
	LocationID        int64     `json:"locationId"`
	LocationName      string    `json:"locationName"`
	ContainerID       *int64    `json:"containerId"`
	DivisionNumber    *int      `json:"divisionNumber"`
	QuantityAvailable int64     `json:"quantityAvailable"`
	PricePerUnit      int64     `json:"pricePerUnit"`
	Notes             *string   `json:"notes"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type PurchaseTransaction struct {
	ID                int64     `json:"id"`
	ForSaleItemID     int64     `json:"forSaleItemId"`
	BuyerUserID       int64     `json:"buyerUserId"`
	BuyerName         string    `json:"buyerName"`
	SellerUserID      int64     `json:"sellerUserId"`
	TypeID            int64     `json:"typeId"`
	TypeName          string    `json:"typeName"`
	LocationID        int64     `json:"locationId"`
	LocationName      string    `json:"locationName"`
	QuantityPurchased int64     `json:"quantityPurchased"`
	PricePerUnit      int64     `json:"pricePerUnit"`
	TotalPrice        int64     `json:"totalPrice"`
	Status            string    `json:"status"`
	ContractKey       *string   `json:"contractKey,omitempty"`
	TransactionNotes  *string   `json:"transactionNotes"`
	PurchasedAt       time.Time `json:"purchasedAt"`
}

type BuyOrder struct {
	ID               int64     `json:"id"`
	BuyerUserID      int64     `json:"buyerUserId"`
	TypeID           int64     `json:"typeId"`
	TypeName         string    `json:"typeName"`
	QuantityDesired  int64     `json:"quantityDesired"`
	MaxPricePerUnit  int64     `json:"maxPricePerUnit"`
	Notes            *string   `json:"notes"`
	IsActive         bool      `json:"isActive"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// Sales Analytics Models

type SalesMetrics struct {
	TotalRevenue      int64          `json:"totalRevenue"`
	TotalTransactions int64          `json:"totalTransactions"`
	TotalQuantitySold int64          `json:"totalQuantitySold"`
	UniqueItemTypes   int64          `json:"uniqueItemTypes"`
	UniqueBuyers      int64          `json:"uniqueBuyers"`
	TimeSeriesData    []TimeSeriesData `json:"timeSeriesData"`
	TopItems          []ItemSalesData  `json:"topItems"`
}

type TimeSeriesData struct {
	Date              string `json:"date"`
	Revenue           int64  `json:"revenue"`
	Transactions      int64  `json:"transactions"`
	QuantitySold      int64  `json:"quantitySold"`
}

type ItemSalesData struct {
	TypeID            int64   `json:"typeId"`
	TypeName          string  `json:"typeName"`
	QuantitySold      int64   `json:"quantitySold"`
	Revenue           int64   `json:"revenue"`
	TransactionCount  int64   `json:"transactionCount"`
	AveragePricePerUnit int64 `json:"averagePricePerUnit"`
}

type BuyerAnalytics struct {
	BuyerUserID       int64   `json:"buyerUserId"`
	BuyerName         string  `json:"buyerName"`
	TotalSpent        int64   `json:"totalSpent"`
	TotalPurchases    int64   `json:"totalPurchases"`
	TotalQuantity     int64   `json:"totalQuantity"`
	FirstPurchaseDate time.Time `json:"firstPurchaseDate"`
	LastPurchaseDate  time.Time `json:"lastPurchaseDate"`
	RepeatCustomer    bool    `json:"repeatCustomer"`
}
