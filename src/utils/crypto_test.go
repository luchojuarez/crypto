package utils

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.BeforeSuite(func() {
	// block all HTTP requests
	//httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
})

var _ = ginkgo.BeforeEach(func() {
	// remove any mocks
	httpmock.Reset()
})

var _ = ginkgo.AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

func TestFoo(t *testing.T) {
	cryptoService := NewCryptoService()
	httpmock.ActivateNonDefault(cryptoService.RestClient.GetClient())
	//httpmock.ActivateNonDefault(resty.New().GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock API response
	responseMap := map[string]interface{}{
		"USD": 6209.29,
		"JPY": 689215.38,
		"EUR": 5830.56,
		"ARS": 551902.3,
	}
	httpmock.RegisterResponder(
		"GET",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, responseMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
	//httpmock.NewStringResponder(200, `{"USD": 6209.29, "JPY": 689215.38, "EUR": 5830.56, "ARS": 551902.3}`))

	call, _ := cryptoService.GetValues()
	expectedARS := float64(551902.3)
	assert.Equal(t, expectedARS, *call.Ars)
	assert.Equal(t, "$551.902,30", call.ArsToString())

}

func TestApiDown(t *testing.T) {
	cryptoService := NewCryptoService()
	httpmock.ActivateNonDefault(cryptoService.RestClient.GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock API response
	httpmock.RegisterResponder(
		"GET",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS",
		httpmock.NewStringResponder(500, `{"USD": 0, "JPY: 689215.38, "EUR": 5830.56, "ARS": 551902.3}`))

	_, err := cryptoService.GetValues()
	assert.NotNil(t, err, fmt.Sprintf("err: '%v'", err))
	assert.Equal(t, "invalid status code: '500'", err.Error())
}

func TestInvalidJSON(t *testing.T) {
	cryptoService := NewCryptoService()
	httpmock.ActivateNonDefault(cryptoService.RestClient.GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock API response
	httpmock.RegisterResponder(
		"GET",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS",
		httpmock.NewStringResponder(200, `{"USD": 0, "JPY: 689215.38, "EUR": 5830.56, "ARS": 551902.3}`))

	_, err := cryptoService.GetValues()
	assert.NotNil(t, err, fmt.Sprintf("err: '%v'", err))
	assert.Equal(t, "invalid character 'E' after object key", err.Error())
}

func TestRestError(t *testing.T) {
	cryptoService := NewCryptoService()
	httpmock.ActivateNonDefault(cryptoService.RestClient.GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// NO mocks, rest error
	_, err := cryptoService.GetValues()
	assert.NotNil(t, err, fmt.Sprintf("err: '%v'", err))
}

func TestSlow(t *testing.T) {

	startMillis := time.Now().UnixNano() / int64(time.Millisecond)

	cryptoService := NewCryptoService()
	httpmock.ActivateNonDefault(cryptoService.RestClient.GetClient())
	//httpmock.ActivateNonDefault(resty.New().GetClient())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// mock API response
	responseMap := map[string]interface{}{
		"USD": 6209.29,
		"JPY": 689215.38,
		"EUR": 5830.56,
		"ARS": 551902.3,
	}
	httpmock.RegisterResponder(
		"GET",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(1 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, responseMap)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
	//httpmock.NewStringResponder(200, `{"USD": 6209.29, "JPY": 689215.38, "EUR": 5830.56, "ARS": 551902.3}`))

	go cryptoService.GetValues()
	// log.Printf("esto tiene '%v'", call)
	// expectedARS := float64(551902.3)
	// assert.Equal(t, expectedARS, *call.Ars)
	// assert.Equal(t, "$551.902,30", call.ArsToString())

	endMillis := time.Now().UnixNano() / int64(time.Millisecond)

	log.Printf("desde el test '%d'", endMillis-startMillis)

}
