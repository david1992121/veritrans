package veritrans

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	assert "github.com/stretchr/testify/require"
)

var accountService *AccountService

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file for testing")
	}
	accountService = NewAccountService(AccountConfig{
		MerchantCCID:     os.Getenv("MERCHANT_CCID"),
		MerchantPassword: os.Getenv("MERCHANT_PASSWORD"),
		ApiURL:           os.Getenv("CONNECTION_API_URL"),
		TxnVersion:       os.Getenv("TXN_VERSION"),
		DummyRequest:     os.Getenv("DUMMY_REQUEST"),
	})
}

func TestAccount(t *testing.T) {
	testAccountID := "ACCOUNT_SERVICE_001"
	accountParam := &AccountParam{
		AccountID: testAccountID,
	}

	// Get Account
	account, err := accountService.GetAccount(accountParam)
	if err == nil {
		// Assert if the account exists
		assert.Equal(t, testAccountID, account.AccountID)
	} else {
		assert.Equal(t, "未登録の会員です。", err.Error())
		account, err := accountService.CreateAccount(accountParam)

		// Create if the account doesn't exist
		assert.Nil(t, err)
		assert.Equal(t, testAccountID, account.AccountID)
	}

	// Remove account
	account, err = accountService.DeleteAccount(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, testAccountID, account.AccountID)

	// Restore account
	account, err = accountService.RestoreAccount(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, testAccountID, account.AccountID)
}
