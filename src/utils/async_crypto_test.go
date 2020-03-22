package utils

import (
	"log"
	"net/http"
	"github.com/luchojuarez/src/models"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFoo2(t *testing.T) {
	// time trace
	defer printTotaTime(time.Now(), "tiempo del lento")

	// news channel
	c := make(chan models.CryptoCompareResponse)
	// service instance
	cryptoService := NewAsyncCriptoService()
	//mock HTTP server
	httpmock.ActivateNonDefault(cryptoService.MainService.RestClient.GetClient())
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	// mock API response
	responseMap := map[string]interface{}{
		"USD": 6209.29,
		"JPY": 689215.38,
		"EUR": 5830.56,
		"ARS": 551902.3,
	}

	// mock response API
	// usea a mock interceptor to sleep request
	httpmock.RegisterResponder(
		"GET",
		"https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR,ARS",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, responseMap)
			time.Sleep(1 * time.Second)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	go cryptoService.GetValues(c)
	call := <-c
	expectedARS := float64(551902.3)
	assert.Equal(t, expectedARS, *call.Ars)
	assert.Equal(t, "$551.902,30", call.ArsToString())

	log.Printf("desde el test '%v', to str: '%s'", call, call.ArsToString())

}

func printTotaTime(startTime time.Time, message string) {
	beginMillis := startTime.UnixNano() / int64(time.Millisecond)
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("%s: %dms", message, endMillis-beginMillis)
}
