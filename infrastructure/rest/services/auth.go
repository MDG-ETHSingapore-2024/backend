package services

func VerifyWalletID(walletId string) bool {
	if walletId == "" {
		return false
	}
	return true
}
