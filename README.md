[![CI](https://github.com/david1992121/veritrans/actions/workflows/main.yml/badge.svg)](https://github.com/david1992121/veritrans/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/david1992121/veritrans/branch/main/graph/badge.svg?token=IJAN6H2PU2)](https://codecov.io/gh/david1992121/veritrans)
![GitHub](https://img.shields.io/github/license/david1992121/veritrans?label=license)
[![Go Report Card](https://goreportcard.com/badge/github.com/david1992121/veritrans)](https://goreportcard.com/report/github.com/david1992121/veritrans)

# Go Veritrans4G

## Overview

- Account and card management (Add/Update/Delete/Restore/Get)
- Payment with the MDK token
- Payment with the registered account

## Getting Started

### MDK Token Payment ###

1. Initialize the MDK service
```
cardService := NewMDKService(MDKConfig{
  APIURL:   "MDK_API_URL",
  APIToken: "MDK_API_TOKEN",
})
```
2. Get MDK token from the card information
```
cardToken, err := cardService.GetCardToken(&ClientCardInfo{
  CardNumber:   "4111111111111111",
  CardExpire:   "12/25",
  SecurityCode: "123",
})
```
3. Initialize the payment service
```
config := ConnectionConfig{
  MerchantCCID:     "MERCHANT_CCID",
  MerchantPassword: "MERCHANT_PASSWORD",
  PaymentAPIURL:    "PAYMENT_API_URL",
  TxnVersion:       "TXN_VERSION",
  DummyRequest:     "DUMMY_REQUEST",
}
paymentService, _ = NewPaymentService(config)
```
4. Pay with the MDK token
```
authorizeParam := Params{
  OrderID:     "MDK Order ID",
  Amount:      "Amount",
  JPO:         "10",
  WithCapture: "true",
  PayNowIDParam: &PayNowIDParam{
    Token: "MDK TOKEN"
  },
}
paymentService.Authorize(&authorizeParam, PaymentServiceType(PayCard))
```

### Account Management ###

1. Initialize the account service
```
accountService = NewAccountService(ConnectionConfig{
  MerchantCCID:    "MERCHANT_CCID",
  MerchantPassword:"MERCHANT_PASSWORD",
  AccountAPIURL:   "ACCOUNT_API_URL",
  TxnVersion:      "TXN_VERSION",
  DummyRequest:    "DUMMY_REQUEST",
})
accountService.CreateAccount(accountParam)
```
2. Create a veritrans account
```
accountParam := &AccountParam{
  AccountID: "Your Account ID",
}
accountService.CreateAccount(accontParam)
```
3. Add a credit card
```
accountParam.CardParam = &CardParam{
  CardNumber:  "Your Card Number",
  CardExpire:  "12/25",
  DefaultCard: "1",
}
accountService.CreateCard(accountParam)
```
4. Pay with the registered account
```
authorizeParam := Params{
  OrderID:     "Account Order ID",
  Amount:      payAmount,
  JPO:         "10",
  WithCapture: "false",
  PayNowIDParam: &PayNowIDParam{
    AccountParam: &AccountParam{
      AccountID: "Your Account ID",
    },
  },
}
paymentService.Authorize(&authorizeParam, PaymentServiceType(PayCard))
```

## Test ##

Create .env file and define the required variables
```
MDK_API_TOKEN=xxxxxxxx-xxxx-xxxx-xxxxxxxxxx
MDK_API_URL=https://api.veritrans.co.jp/4gtoken
MERCHANT_CCID=A100000000000001069713cc
MERCHANT_PASSWORD=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
DUMMY_REQUEST=0
TXN_VERSION=2.0.0
ACCOUNT_API_URL=https://api.veritrans.co.jp:443/paynowid/v1
PAYMENT_API_URL=https://api.veritrans.co.jp:443/paynow/v2
SEARCH_API_URL=https://api.veritrans.co.jp:443/paynow-search/v2
```
Run the test using the following command
```console
$ go test -v
```