package veritrans

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AccountService struct {
	Config AccountConfig
}

// Configuration of the account service
// ApiURL is the account management api endpoint (https://api.veritrans.co.jp:443/paynowid/v1/)
// TxnVersion is the version of the veritrans api (2.0.0)
// DummyRequest is the flag indicating whether the request is dummy or live
type AccountConfig struct {
	MerchantCCID     string
	MerchantPassword string
	ApiURL           string
	TxnVersion       string
	DummyRequest     string
}

func NewAccountService(config AccountConfig) *AccountService {
	return &AccountService{Config: config}
}

func (acc AccountService) ExecuteAccountRequest(requestURL string, connectionParam *ConnectionParam) (*AccountResponse, error) {
	paramByte, err := json.Marshal(connectionParam)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.ParseRequestURI(requestURL)
	if err == nil {
		httpClient := &http.Client{}
		body := bytes.NewBuffer(paramByte)
		req, err := http.NewRequest("POST", parsedURL.String(), body)
		req.Header.Set("Content-Type", "application/json")

		if err == nil {
			res, err := httpClient.Do(req)
			if err == nil {
				body, err := ioutil.ReadAll(res.Body)
				if err == nil {
					var accountRes AccountResponse
					err = json.Unmarshal(body, &accountRes)
					return &accountRes, err
				}
			}
		}
	}
	return nil, err
}

// Handler to make hash data of params
func (acc AccountService) SetHash(connectionParam *ConnectionParam) error {
	paramJSON, err := json.Marshal(connectionParam.Params)
	if err != nil {
		return err
	}

	hash := []byte(fmt.Sprintf("%s%s%s", acc.Config.MerchantCCID, paramJSON, acc.Config.MerchantPassword))

	sha := sha256.New()
	sha.Write(hash)
	connectionParam.AuthHash = fmt.Sprintf("%x", sha.Sum(nil))
	return nil
}

// Get connection paramater from account parameter
func (acc AccountService) GetConnectionParam(accountParam *AccountParam) (*ConnectionParam, error) {
	payNowIdParam := &PayNowIDParam{
		AccountParam: accountParam,
	}
	payNowIdParam.Default()

	connectionParam := &ConnectionParam{
		Params: Params{
			PayNowIDParam: *payNowIdParam,
			TxnVersion:    acc.Config.TxnVersion,
			DummyRequest:  acc.Config.DummyRequest,
			MerchantCCID:  acc.Config.MerchantCCID,
		},
		AuthHash: "",
	}

	if err := acc.SetHash(connectionParam); err != nil {
		return nil, err
	}
	return connectionParam, nil
}

// Execute Account CRUD
func (acc AccountService) ExecuteAccountProcess(serviceType ServiceType, mode ManagementMode, accountParam *AccountParam) (*Account, error) {
	connectionParam, err := acc.GetConnectionParam(accountParam)
	if err == nil {
		accountRes, err := acc.ExecuteAccountRequest(
			fmt.Sprintf("%s/%s/%s", acc.Config.ApiURL, managementModes[mode], serviceTypes[serviceType]), connectionParam)
		if err == nil {
			if accountRes.Result.MStatus == "success" {
				return &accountRes.PayNowIDResponse.Account, nil
			}

			return nil, errors.New(accountRes.Result.MErrorMsg)
		}
	}
	return nil, err
}

// Create a veritrans account
func (acc AccountService) CreateAccount(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(AccountType, MethodAdd, accountParam)
}

// Remove a veritrans account
func (acc AccountService) DeleteAccount(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(AccountType, MethodDelete, accountParam)
}

// Get a veritrans account
func (acc AccountService) GetAccount(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(AccountType, MethodGet, accountParam)
}

// Get a veritrans account
func (acc AccountService) RestoreAccount(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(AccountType, MethodRestore, accountParam)
}

// Create a card
func (acc AccountService) CreateCard(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(CardType, MethodAdd, accountParam)
}

// Remove a card
func (acc AccountService) DeleteCard(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(CardType, MethodDelete, accountParam)
}

// Update a card information
func (acc AccountService) UpdateCard(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(CardType, MethodUpdate, accountParam)
}

// Get a veritrans account
func (acc AccountService) GetCard(accountParam *AccountParam) (*Account, error) {
	return acc.ExecuteAccountProcess(CardType, MethodGet, accountParam)
}
