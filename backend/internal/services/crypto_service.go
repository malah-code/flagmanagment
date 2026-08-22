package services

import (
	"github.com/flagmanagment/backend/internal/crypto"
)

type CryptoService = crypto.CryptoService

func NewCryptoService() CryptoService {
	return crypto.NewCryptoService()
}
