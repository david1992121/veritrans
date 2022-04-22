package veritrans

// Default interface fills default values
type Default interface {
	Default()
}

// AccountBasicParam represents the "accountBasicParam" of the request.
type AccountBasicParam struct {
	CreateDate      string `json:"createDate"`
	DeleteDate      string `json:"deleteDate"`
	ForceDeleteDate string `json:"forceDeleteDate"`
}

// CardParam is represents the "cardParam" of the request.
type CardParam struct {
	CardID        string `json:"cardId"`
	DefaultCard   string `json:"defaultCard,omitempty"`
	DefaultCardID string `json:"defaultCardId,omitempty"`
	CardNumber    string `json:"cardNumber"`
	CardExpire    string `json:"cardExpire"`
	Token         string `json:"token,omitempty"`
}

// RecurringChargeParm represents the "recurringChargeParam" of the request.
type RecurringChargeParam struct {
	GroupID       string `json:"groupId"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	FinalCharge   string `json:"finalCharge,omitempty"`
	OneTimeAmount string `json:"oneTimeAmount"`
	Amount        string `json:"amount"`
}

// AccountParam represents the "accountParam" of the request.
type AccountParam struct {
	AccountID            string                `json:"accountId"`
	AccountBasicParam    *AccountBasicParam    `json:"accountBasicParam,omitempty"`
	CardParm             *CardParam            `json:"cardParam,omitempty"`
	RecurringChargeParam *RecurringChargeParam `json:"recurringChargeParam,omitempty"`
}

// PayNowIDParm represents the "payNowIdParam" of the request.
type PayNowIDParam struct {
	Token        string        `json:"token,omitempty"`
	AccountParam *AccountParam `json:"accountParam,omitempty"`
	Memo         string        `json:"memo1,omitempty"`
	FreeKey      string        `json:"freeKey,omitempty"`
}

// Params represents the "params" of the request.
type Params struct {
	OrderID       string        `json:"orderId,omitempty"`
	Amount        string        `json:"amount,omitempty"`
	JPO           string        `json:"jpo,omitempty"`
	WithCapture   string        `json:"withCapture,omitempty"`
	PayNowIDParam PayNowIDParam `json:"payNowIdParam"`
	TxnVersion    string        `json:"txnVersion"`
	DummyRequest  string        `json:"dummyRequest"`
	MerchantCCID  string        `json:"merchantCcid"`
}

// ConnectionParam represents the request parameter.
type ConnectionParam struct {
	Params   Params `json:"params"`
	AuthHash string `json:"authHash"`
}

// implementations of the Default interface
func (payParam *PayNowIDParam) Default() {
	if payParam.Memo == "" {
		payParam.Memo = "memo"
	}
	if payParam.FreeKey == "" {
		payParam.FreeKey = "freekey"
	}
}

func (accountBasicParam *AccountBasicParam) Default() {
	if accountBasicParam.ForceDeleteDate == "" {
		accountBasicParam.ForceDeleteDate = "0"
	}
}

func (recurringChargeParam *RecurringChargeParam) Default() {
	if recurringChargeParam.FinalCharge == "" {
		recurringChargeParam.FinalCharge = "0"
	}
}

// response types
type Result struct {
	VResultCode string `json:"vResultCode"`
	MStatus     string `json:"mstatus"`
	MErrorMsg   string `json:"merrMsg"`
}

type CardInfo struct {
	CardExpire  string `json:"cardExpire"`
	CardID      string `json:"cardId"`
	CardNumber  string `json:"cardNumber"`
	DefaultCard string `json:"defaultCard"`
}

type Account struct {
	AccountID string     `json:"accountId"`
	CardInfo  []CardInfo `json:"cardInfo"`
}

type PayNowIDResponse struct {
	Account Account `json:"account"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

type AccountResponse struct {
	PayNowIDResponse *PayNowIDResponse `json:"payNowIdResponse,omitempty"`
	Result           Result            `json:"result"`
}
