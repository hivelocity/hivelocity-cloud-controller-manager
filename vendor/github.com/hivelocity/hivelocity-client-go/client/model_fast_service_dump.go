/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type FastServiceDump struct {
	ServiceId                int32               `json:"serviceId,omitempty"`
	Created                  int32               `json:"created,omitempty"`
	BilledPrice              float32             `json:"billedPrice,omitempty"`
	BilledPricePerPeriod     float32             `json:"billedPricePerPeriod,omitempty"`
	ServiceCost              float32             `json:"serviceCost,omitempty"`
	LastRenew                int32               `json:"lastRenew,omitempty"`
	RenewDate                int32               `json:"renewDate,omitempty"`
	Quantity                 float32             `json:"quantity,omitempty"`
	OrderId                  int32               `json:"orderId,omitempty"`
	Status                   string              `json:"status,omitempty"`
	Period                   string              `json:"period,omitempty"`
	Discount                 float32             `json:"discount,omitempty"`
	DiscountType             string              `json:"discountType,omitempty"`
	DiscountedCost           float32             `json:"discountedCost,omitempty"`
	ServiceDiscount          float32             `json:"serviceDiscount,omitempty"`
	ServiceDiscountPerPeriod float32             `json:"serviceDiscountPerPeriod,omitempty"`
	IpAddress                string              `json:"ipAddress,omitempty"`
	CancelAfter              int32               `json:"cancelAfter,omitempty"`
	StartDate                int32               `json:"startDate,omitempty"`
	EndDate                  int32               `json:"endDate,omitempty"`
	ServiceOptions           []ServiceOptionData `json:"serviceOptions,omitempty"`
	Usage                    interface{}         `json:"usage,omitempty"`
	ServiceDevices           []interface{}       `json:"serviceDevices,omitempty"`
	ChildServices            []interface{}       `json:"childServices,omitempty"`
	ProductId                int32               `json:"productId,omitempty"`
	ProductName              string              `json:"productName,omitempty"`
	Reseller                 string              `json:"reseller,omitempty"`
	ServiceType              string              `json:"serviceType,omitempty"`
	ContractTerm             int32               `json:"contractTerm,omitempty"`
	BillingInfoId            int32               `json:"billingInfoId,omitempty"`
	AutoBill                 bool                `json:"autoBill,omitempty"`
	ColocationCharge         float32             `json:"colocationCharge,omitempty"`
	ResellerBmaasCharge      float32             `json:"resellerBmaasCharge,omitempty"`
	FacilityName             string              `json:"facilityName,omitempty"`
	// The service type code. The list of service types can be accessed on https://core.hivelocity.net/api/v2/service/types .
	TypeCode string `json:"typeCode,omitempty"`
}
