package util

// List for all supported currencies
var USD = "USD"
var EUR = "EUR"
var BRL = "BRL"

var CURRENCY_LIST = [3]string{USD, EUR, BRL}

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	for _, cur := range CURRENCY_LIST {
		if cur == currency {
			return true
		}
	}
	return false
}
