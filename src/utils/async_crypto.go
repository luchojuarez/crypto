package utils

import "github.com/luchojuarez/crypto/src/models"

type AsyncCriptoService struct {
	MainService CryptoService
}

func NewAsyncCriptoService() AsyncCriptoService {
	return AsyncCriptoService{
		MainService: NewCryptoService(),
	}
}

func (self AsyncCriptoService) GetValues(c chan models.CryptoCompareResponse) {

	toReturn, _ := self.MainService.GetValues()
	c <- *toReturn
}
