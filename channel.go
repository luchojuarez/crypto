package main

import (
	"log"

	"github.com/luchojuarez/crypto/src/utils"
)

func main() {

	cryptoResponse, _ := utils.NewCryptoService().GetValues()

	log.Printf("BTC to ARS value '%s' response time: %d", cryptoResponse.ArsToString(), cryptoResponse.ResponseTime)

}
