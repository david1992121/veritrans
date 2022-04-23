package veritrans

import (
	"fmt"
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
	testAccountID := "TEST_ACCOUNT_01"
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
	fmt.Println("Create Account Passed")

	// Remove account
	account, err = accountService.DeleteAccount(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, testAccountID, account.AccountID)
	fmt.Println("Remove Account Passed")

	// Restore account
	account, err = accountService.RestoreAccount(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, testAccountID, account.AccountID)
	fmt.Println("Restore Account Passed")
}

func TestCard(t *testing.T) {
	testAccountID := "TEST_ACCOUNT_02"
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

	// Add Card
	testCardNumber := "4111111111111111"
	expectedCardNumber := "411111********11"
	expiredAt := getAfterOneMonth()
	accountParam.CardParam = &CardParam{
		CardNumber:  testCardNumber,
		CardExpire:  getAfterOneMonth(),
		DefaultCard: "1",
	}

	account, err = accountService.CreateCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(account.CardInfo))
	fmt.Println("Add Card Passed")

	// Get Card
	accountParam.CardParam = nil
	account, err = accountService.GetCard(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(account.CardInfo))
	assert.Equal(t, expectedCardNumber, account.CardInfo[0].CardNumber)
	assert.Equal(t, expiredAt, account.CardInfo[0].CardExpire)
	fmt.Println("Get Card Passed")
	cardID := account.CardInfo[0].CardID

	// Update Card
	newExpiredAt := getAfterOneYear()
	accountParam.CardParam = &CardParam{
		CardID:     cardID,
		CardExpire: newExpiredAt,
	}
	account, err = accountService.UpdateCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(account.CardInfo))
	assert.Equal(t, newExpiredAt, account.CardInfo[0].CardExpire)
	fmt.Println("Update Card Passed")

	// Remove Card
	accountParam.CardParam = &CardParam{
		CardID: cardID,
	}
	account, err = accountService.DeleteCard(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(account.CardInfo))
	fmt.Println("Remove Card Passed")
}
