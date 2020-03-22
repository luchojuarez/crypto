package models

import (
	"github.com/leekchan/accounting"
)

type CryptoCompareResponse struct {
	Ars                 *float64 `json:"ARS"`
	Usd                 *float64 `json:"USD"`
	ResponseTime        int64
	ResponseStatusConde int
}

func (this CryptoCompareResponse) ArsToString() string {
	if this.Ars == nil {
		return ""
	}
	ac := accounting.Accounting{Symbol: "$", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(*this.Ars)
}
