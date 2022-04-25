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
		for _, envItem := range EnvVariables {
			if os.Getenv(envItem) == "" {
				log.Fatal("No env file for testing")
			}
		}
	}

	accountService = NewAccountService(ConnectionConfig{
		MerchantCCID:     os.Getenv("MERCHANT_CCID"),
		MerchantPassword: os.Getenv("MERCHANT_PASSWORD"),
		AccountAPIURL:    os.Getenv("ACCOUNT_API_URL"),
		TxnVersion:       os.Getenv("TXN_VERSION"),
		DummyRequest:     os.Getenv("DUMMY_REQUEST"),
	})
}

func TestAccount(t *testing.T) {
	testAccountID := "TEST_ACCOUNT_1"
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
	testAccountID := "TEST_ACCOUNT_2"
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
	firstTestCardNumber := "4111111111111111"
	firstExpectedCardNumber := "411111********11"
	expiredAt := GetAfterOneMonth()
	accountParam.CardParam = &CardParam{
		CardNumber:  firstTestCardNumber,
		CardExpire:  expiredAt,
		DefaultCard: "1",
	}

	account, err = accountService.CreateCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(account.CardInfo))
	fmt.Println("Add The First Card Passed")

	// Get Cards
	accountParam.CardParam = nil
	account, err = accountService.GetCard(accountParam)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(account.CardInfo))
	assert.Equal(t, firstExpectedCardNumber, account.CardInfo[0].CardNumber)
	assert.Equal(t, "1", account.CardInfo[0].DefaultCard)
	assert.Equal(t, expiredAt, account.CardInfo[0].CardExpire)
	fmt.Println("Get Card Passed")
	firstCardID := account.CardInfo[0].CardID

	// Add another card
	secondTestCardNumber := "5555555555554444"
	secondExpectedCardNumber := "555555********44"
	accountParam.CardParam = &CardParam{
		CardNumber:  secondTestCardNumber,
		CardExpire:  expiredAt,
		DefaultCard: "0",
	}

	account, err = accountService.CreateCard(accountParam)
	assert.Nil(t, err)
	secondCardID := account.CardInfo[0].CardID
	fmt.Println("Add The Second Card Passed")

	// Get Cards
	accountParam.CardParam = nil
	account, err = accountService.GetCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(account.CardInfo))
	assert.Equal(t, secondExpectedCardNumber, account.CardInfo[1].CardNumber)
	assert.Equal(t, "0", account.CardInfo[1].DefaultCard)

	// Update Card
	newExpiredAt := GetAfterOneYear()
	accountParam.CardParam = &CardParam{
		CardID:      secondCardID,
		DefaultCard: "1",
		CardExpire:  newExpiredAt,
	}
	account, err = accountService.UpdateCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, secondExpectedCardNumber, account.CardInfo[0].CardNumber)
	assert.Equal(t, newExpiredAt, account.CardInfo[0].CardExpire)
	assert.Equal(t, "1", account.CardInfo[0].DefaultCard)

	// Remove Two Cards
	accountParam.CardParam = &CardParam{
		CardID: firstCardID,
	}
	_, err = accountService.DeleteCard(accountParam)
	assert.Nil(t, err)

	accountParam.CardParam = &CardParam{
		CardID: secondCardID,
	}
	_, err = accountService.DeleteCard(accountParam)
	assert.Nil(t, err)

	// Get Cards
	accountParam.CardParam = nil
	account, err = accountService.GetCard(accountParam)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(account.CardInfo))
	fmt.Println("Remove Card Passed")
}
