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
// URL is the account management api endpoint (https://api.veritrans.co.jp:443/paynowid/v1/)
type AccountConfig struct {
	MerchantCCID     string
	MerchantPassword string
	URL              string
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
					var accountRes *AccountResponse
					err = json.Unmarshal(body, accountRes)
					return accountRes, err
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

// Create a veritrans account
func (acc AccountService) CreateAccount(accountParam *AccountParam) (*Account, error) {
	payNowIdParam := &PayNowIDParam{
		AccountParam: accountParam,
	}
	payNowIdParam.Default()

	connectionParam := ConnectionParam{
		Params: Params{
			PayNowIDParam: *payNowIdParam,
			TxnVersion:    acc.Config.TxnVersion,
			DummyRequest:  acc.Config.DummyRequest,
		},
		AuthHash: "",
	}

	acc.SetHash(&connectionParam)

	accountRes, err := acc.ExecuteAccountRequest(acc.Config.URL+"Add/account", &connectionParam)
	if err != nil {
		return nil, err
	}

	if accountRes.Result.MStatus == "success" {
		return &accountRes.PayNowIDResponse.Account, nil
	}

	return nil, errors.New(accountRes.Result.MErrorMsg)
}
