package services

// Method to verify access tokens sent with API calls
// TODO: integrate with the web3 track we will use
func VerifyToken(token string) bool {
	// TODO: put some real logic
	return token == "test_token"
}