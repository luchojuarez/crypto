package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luchojuarez/crypto/src/models"

	"github.com/go-resty/resty/v2"
)

const cryptocompareBaseURL = "https://min-api.cryptocompare.com"

type CryptoService struct {
	RestClient *resty.Client
}

// NewCryptoService creates a CryptoService
func NewCryptoService() CryptoService {
	return CryptoService{
		RestClient: resty.New(),
	}
}

func (self CryptoService) GetValues() (*models.CryptoCompareResponse, error) {
	startMillis := time.Now().UnixNano() / int64(time.Millisecond)
	SC, body, err := self.simpleGet()
	if err != nil {
		return nil, err
	}
	var toReturn models.CryptoCompareResponse
	if err = json.Unmarshal([]byte(body), &toReturn); err != nil {
		return nil, err
	}
	toReturn.ResponseStatusConde = SC
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)

	toReturn.ResponseTime = (endMillis - startMillis)
	return &toReturn, nil
}

// simple wrapper to resty
func (self CryptoService) simpleGet() (int, string, error) {

	response, err := self.RestClient.
		R().
		// SetQueryParams(map[string]string{
		// 	"fsym":  "BTC",               //from currency id
		// 	"tsyms": "USD,JPY,EUR,ARS"}). // to currency id
		//Get(cryptocompareBaseURL + "/data/price")
		Get(cryptocompareBaseURL + "/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS")
	if err != nil {
		return 0, "", err
	}
	if response.StatusCode() != http.StatusOK {
		return response.StatusCode(), "", fmt.Errorf("invalid status code: '%d'", response.StatusCode())
	}

	return response.StatusCode(), fmt.Sprintf("%s", response), nil
}
