package data

import "time"

type ImgCdrOperatorRoutingData struct {
	RoutingPrefix              string     `json:"RoutingPrefix"`
	BaseOperator               *string    `json:"BaseOperator"`
	BaseClass1                 *int       `json:"BaseClass1"`
	BaseClass2                 *int       `json:"BaseClass2"`
	BaseClass3                 *int       `json:"BaseClass3"`
	BaseClass4                 *int       `json:"BaseClass4"`
	BaseClassCreationDate      *time.Time `json:"BaseClassCreationDate"`
	CurrentClass1              *int       `json:"CurrentClass1"`
	CurrentClass2              *int       `json:"CurrentClass2"`
	CurrentClass3              *int       `json:"CurrentClass3"`
	CurrentClass4              *int       `json:"CurrentClass4"`
	CurrentClassLastUpdateDate *time.Time `json:"CurrentClassLastUpdateDate"`
	CurrentClassRemarks        *string    `json:"CurrentClassRemarks"`
	RoutingManualOverride      *string    `json:"RoutingManualOverride"`
	OverrideCreationDate       *time.Time `json:"OverrideCreationDate"`
	OverrideCreatedBy          *string    `json:"OverrideCreatedBy"`
	OverrideReason             *string    `json:"OverrideReason"`
	BaseClassLastUpdateDate    *time.Time `json:"BaseClassLastUpdateDate"`
}
