package types

// CardRequest represents the request params of the MDK token api
type CardRequest struct {
	CardNumber     string `json:"card_number"`
	CardExpire     string `json:"card_expire"`
	SecurityCode   string `json:"security_code"`
	CardHolderName string `json:"cardholder_name,omitempty"`
	TokenAPIKey    string `json:"token_api_key"`
	Lang           string `json:"lang"`
}

// CardResponse represents the response of the MDK token api
type CardResponse struct {
	Token           string `json:"token"`
	TokenExpireDate string `json:"token_expire_date"`
	ReqCardNumber   string `json:"req_card_number"`
	Status          string `json:"string"`
	Code            string `json:"code"`
	Message         string `json:"message"`
}
