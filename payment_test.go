package veritrans

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	assert "github.com/stretchr/testify/require"
)

var paymentService *PaymentService

func init() {
	if err := godotenv.Load(); err != nil {
		for _, envItem := range EnvVariables {
			if os.Getenv(envItem) == "" {
				log.Fatal("No env file for testing")
			}
		}
	}

	config := ConnectionConfig{
		MerchantCCID:     os.Getenv("MERCHANT_CCID"),
		MerchantPassword: os.Getenv("MERCHANT_PASSWORD"),
		AccountAPIURL:    os.Getenv("ACCOUNT_API_URL"),
		PaymentAPIURL:    os.Getenv("PAYMENT_API_URL"),
		SearchAPIURL:     os.Getenv("SEARCH_API_URL"),
		TxnVersion:       os.Getenv("TXN_VERSION"),
		DummyRequest:     os.Getenv("DUMMY_REQUEST"),
	}
	paymentService, _ = NewPaymentService(config)
	accountService = NewAccountService(config)
}

func TestPayment(t *testing.T) {
	testAccountID := "PAYMENT_ACCOUNT_01"
	accountParam := &AccountParam{
		AccountID: testAccountID,
	}

	// Create Account
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

	account, err = accountService.GetCard(accountParam)
	assert.Nil(t, err)
	if len(account.CardInfo) == 0 {
		// Add Card
		firstTestCardNumber := "4111111111111111"
		expiredAt := GetAfterOneYear()
		accountParam.CardParam = &CardParam{
			CardNumber:  firstTestCardNumber,
			CardExpire:  expiredAt,
			DefaultCard: "1",
		}
		_, err = accountService.CreateCard(accountParam)
		assert.Nil(t, err)
	}

	// Authorize without capture
	testOrderID, err := findNewOrderID()
	assert.Nil(t, err)
	fmt.Printf("Found New Order ID: %s\n", testOrderID)

	payAmount := "100"
	authorizeParam := Params{
		OrderID:     testOrderID,
		Amount:      payAmount,
		JPO:         "10",
		WithCapture: "false",
		PayNowIDParam: &PayNowIDParam{
			AccountParam: &AccountParam{
				AccountID: testAccountID,
			},
		},
	}
	_, err = paymentService.Authorize(&authorizeParam, PaymentServiceType(PayCard))
	assert.Nil(t, err)

	searchParam := Params{
		ContainDummyFlag: "1",
		ServiceTypeCd:    []string{"card"},
		NewerFlag:        "true",
		SearchParam: &SearchParam{
			Common: OrderParam{
				OrderID: testOrderID,
			},
		},
	}
	result, err := paymentService.Search(&searchParam, PaymentServiceType(Search))
	assert.Nil(t, err)
	assert.Equal(t, "success", result.MStatus)
	assert.NotNil(t, result.OrderInfos)
	assert.GreaterOrEqual(t, len(result.OrderInfos.OrderInfo), 1)
	firstOrderInfo := result.OrderInfos.OrderInfo[0]
	assert.Equal(t, firstOrderInfo.AccountID, testAccountID)
	assert.Equal(t, firstOrderInfo.LastSuccessTxnType, "Authorize")
	assert.NotNil(t, firstOrderInfo.TransactionInfos)
	assert.Equal(t, len(firstOrderInfo.TransactionInfos.TransactionInfo), 1)
	transactionInfo := firstOrderInfo.TransactionInfos.TransactionInfo[0]
	assert.Equal(t, transactionInfo.Amount, payAmount)
	assert.Equal(t, transactionInfo.Command, "Authorize")
	assert.Equal(t, transactionInfo.MStatus, "success")
	assert.Equal(t, transactionInfo.ProperInfo.ReqWithCapture, "false")
	fmt.Println("Authorize Without Capture Passed")

	// Cancel
	cancelParam := Params{
		OrderID: testOrderID,
		Amount:  payAmount,
	}
	_, err = paymentService.Cancel(&cancelParam, PaymentServiceType(PayCard))
	assert.Nil(t, err)

	result, err = paymentService.Search(&searchParam, PaymentServiceType(Search))
	assert.Nil(t, err)
	assert.Equal(t, "success", result.MStatus)
	assert.NotNil(t, result.OrderInfos)
	assert.GreaterOrEqual(t, len(result.OrderInfos.OrderInfo), 1)
	firstOrderInfo = result.OrderInfos.OrderInfo[0]
	assert.Equal(t, firstOrderInfo.AccountID, testAccountID)
	assert.Equal(t, firstOrderInfo.LastSuccessTxnType, "Cancel")
	assert.NotNil(t, firstOrderInfo.TransactionInfos)
	assert.Equal(t, len(firstOrderInfo.TransactionInfos.TransactionInfo), 1)
	transactionInfo = firstOrderInfo.TransactionInfos.TransactionInfo[0]
	assert.Equal(t, transactionInfo.Amount, payAmount)
	assert.Equal(t, transactionInfo.Command, "Cancel")
	assert.Equal(t, transactionInfo.MStatus, "success")
	fmt.Println("Cancel Passed")

	// Authorize with capture
	testOrderID, err = findNewOrderID()
	assert.Nil(t, err)
	fmt.Printf("Found New Order ID: %s\n", testOrderID)

	authorizeParam.WithCapture = "true"
	authorizeParam.OrderID = testOrderID
	_, err = paymentService.Authorize(&authorizeParam, PaymentServiceType(PayCard))
	assert.Nil(t, err)

	searchParam.SearchParam.Common.OrderID = testOrderID
	result, err = paymentService.Search(&searchParam, PaymentServiceType(Search))
	assert.Nil(t, err)
	assert.Equal(t, "success", result.MStatus)
	assert.NotNil(t, result.OrderInfos)
	assert.GreaterOrEqual(t, len(result.OrderInfos.OrderInfo), 1)
	firstOrderInfo = result.OrderInfos.OrderInfo[0]
	assert.Equal(t, firstOrderInfo.AccountID, testAccountID)
	assert.Equal(t, firstOrderInfo.LastSuccessTxnType, "Authorize")
	assert.NotNil(t, firstOrderInfo.TransactionInfos)
	assert.Equal(t, len(firstOrderInfo.TransactionInfos.TransactionInfo), 1)
	transactionInfo = firstOrderInfo.TransactionInfos.TransactionInfo[0]
	assert.Equal(t, transactionInfo.Amount, payAmount)
	assert.Equal(t, transactionInfo.Command, "Authorize")
	assert.Equal(t, transactionInfo.MStatus, "success")
	assert.Equal(t, transactionInfo.ProperInfo.ReqWithCapture, "true")
	fmt.Println("Authorize With Capture Passed")
}

func findNewOrderID() (string, error) {
	testOrderID := ""
	for {
		randomID := GetRandomID(8)
		testOrderID = fmt.Sprintf("PAYMENT_TEST_ORDER_%d", randomID)
		searchParam := Params{
			ContainDummyFlag: "1",
			ServiceTypeCd:    []string{"card"},
			NewerFlag:        "true",
			SearchParam: &SearchParam{
				Common: OrderParam{
					OrderID: testOrderID,
				},
			},
		}
		result, err := paymentService.Search(&searchParam, PaymentServiceType(Search))
		if err != nil {
			return "", err
		}
		if len(result.OrderInfos.OrderInfo) == 0 {
			break
		}
	}
	return testOrderID, nil
}
