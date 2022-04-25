package veritrans

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// CardRequest is a struct of the request parameters of the veritrans MDK token api
type CardRequest struct {
	CardNumber     string `json:"card_number"`
	CardExpire     string `json:"card_expire"`
	SecurityCode   string `json:"security_code"`
	CardHolderName string `json:"cardholder_name,omitempty"`
	TokenAPIKey    string `json:"token_api_key"`
	Lang           string `json:"lang"`
}

// CardResponse represents a response of the veritrans MDK token api
type CardResponse struct {
	Token           string `json:"token"`
	TokenExpireDate string `json:"token_expire_date"`
	ReqCardNumber   string `json:"req_card_number"`
	Status          string `json:"status"`
	Code            string `json:"code"`
	Message         string `json:"message"`
}

// ClientCardInfo indicates the card Information of the client
type ClientCardInfo struct {
	CardNumber     string `json:"card_number"`
	CardExpire     string `json:"card_expire"`
	SecurityCode   string `json:"security_code"`
	CardHolderName string `json:"cardholder_name,omitempty"`
}

// MDKConfig is a configuration of the MDK service
type MDKConfig struct {
	APIURL   string
	APIToken string
}

// MDKService handles the several veritrans APIs for MDK payment
type MDKService struct {
	Config MDKConfig
}

// NewMDKService initializes a MDK service
func NewMDKService(config MDKConfig) *MDKService {
	newService := &MDKService{Config: config}
	return newService
}

// ExecuteCardRequest process the requests
func (mdk *MDKService) ExecuteCardRequest(cardRequest *CardRequest) (*CardResponse, error) {
	cardReqJSON, err := json.Marshal(cardRequest)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.ParseRequestURI(mdk.Config.APIURL)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	body := bytes.NewBuffer(cardReqJSON)
	req, err := http.NewRequest("POST", parsedURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cardResponse CardResponse
	err = json.Unmarshal(resBody, &cardResponse)
	if err != nil {
		return nil, err
	}
	return &cardResponse, nil
}

// GetCardToken gets a card token
func (mdk *MDKService) GetCardToken(cardInfo *ClientCardInfo) (string, error) {
	if cardInfo == nil {
		return "", errors.New("no card information")
	}

	cardRequest := CardRequest{
		CardNumber:     cardInfo.CardNumber,
		CardExpire:     cardInfo.CardExpire,
		CardHolderName: cardInfo.CardHolderName,
		SecurityCode:   cardInfo.SecurityCode,
		TokenAPIKey:    mdk.Config.APIToken,
		Lang:           "ja",
	}

	cardResponse, err := mdk.ExecuteCardRequest(&cardRequest)
	if err != nil {
		return "", err
	}

	if cardResponse.Status == "success" {
		return cardResponse.Token, nil
	}
	return "", errors.New(cardResponse.Message)
}
