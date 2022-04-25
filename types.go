package veritrans

// Configuration of veritrans connection
// AccountApiURL is the account management api endpoint (https://api.veritrans.co.jp:443/paynowid/v1/)
// PaymentApiURL is the payment api endpoint (https://api.veritrans.co.jp:443/paynow/v2)
// TxnVersion is the version of the veritrans api (2.0.0)
// DummyRequest is the flag indicating whether the request is dummy or live
type ConnectionConfig struct {
	MerchantCCID     string
	MerchantPassword string
	AccountApiURL    string
	PaymentApiURL    string
	SearchApiURL     string
	TxnVersion       string
	DummyRequest     string
}

// Default interface fills default values
type Default interface {
	Default()
}

// AccountBasicParam represents the "accountBasicParam" of the request.
type AccountBasicParam struct {
	CreateDate      string `json:"createDate,omitempty"`
	DeleteDate      string `json:"deleteDate,omitempty"`
	ForceDeleteDate string `json:"forceDeleteDate"`
}

// CardParam is represents the "cardParam" of the request.
type CardParam struct {
	CardID        string `json:"cardId,omitempty"`
	DefaultCard   string `json:"defaultCard,omitempty"`
	DefaultCardID string `json:"defaultCardId,omitempty"`
	CardNumber    string `json:"cardNumber,omitempty"`
	CardExpire    string `json:"cardExpire,omitempty"`
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
	CardParam            *CardParam            `json:"cardParam,omitempty"`
	RecurringChargeParam *RecurringChargeParam `json:"recurringChargeParam,omitempty"`
}

// PayNowIDParm represents the "payNowIdParam" of the request.
type PayNowIDParam struct {
	Token        string        `json:"token,omitempty"`
	AccountParam *AccountParam `json:"accountParam,omitempty"`
	Memo         string        `json:"memo1,omitempty"`
	FreeKey      string        `json:"freeKey,omitempty"`
}

// OrderParam
type OrderParam struct {
	OrderID string `json:"orderId"`
}

// SearchParam represents the "searchParameters" of the request.
type SearchParam struct {
	Common OrderParam `json:"common"`
}

// Params represents the "params" of the request.
type Params struct {
	OrderID          string         `json:"orderId,omitempty"`
	Amount           string         `json:"amount,omitempty"`
	JPO              string         `json:"jpo,omitempty"`
	WithCapture      string         `json:"withCapture,omitempty"`
	PayNowIDParam    *PayNowIDParam `json:"payNowIdParam,omitempty"`
	ContainDummyFlag string         `json:"containDummyFlag,omitempty"`
	ServiceTypeCd    []string       `json:"serviceTypeCd,omitempty"`
	NewerFlag        string         `json:"newerFlag"`
	SearchParam      *SearchParam   `json:"searchParameters,omitempty"`
	TxnVersion       string         `json:"txnVersion"`
	DummyRequest     string         `json:"dummyRequest"`
	MerchantCCID     string         `json:"merchantCcid"`
}

// ConnectionParam represents the request parameter.
type ConnectionParam struct {
	Params   Params `json:"params"`
	AuthHash string `json:"authHash"`
}

// Account Management modes
type AccountManagementMode int32

const (
	MethodAdd AccountManagementMode = iota
	MethodUpdate
	MethodDelete
	MethodRestore
	MethodGet
)

var AccountManagementModes = []string{"Add", "Update", "Delete", "Restore", "Get"}

// Account Service Type
type AccountServiceType int32

const (
	AccountType AccountServiceType = iota
	CardType
)

var AccountServiceTypes = []string{"account", "cardinfo"}

// Payment modes
type PaymentManagementMode int32

const (
	MethodAuthorize PaymentManagementMode = iota
	MethodReAuthorize
	MethodCapture
	MethodCancel
	MethodSearch
)

var PaymentManagementModes = []string{"Authorize", "ReAuthorize", "Capture", "Cancel", "Search"}

// Payment Service Type
type PaymentServiceType int32

const (
	PayCard PaymentServiceType = iota
	MPI
	CVS
	EM
	Bank
	UPop
	Paypal
	Saison
	Alipay
	Carrier
	Search
)

var PaymentServiceTypes = []string{"card", "mpi", "cvs", "em", "bank", "upop", "paypal", "saison", "alipay", "carrier", "search"}

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
	VResultCode string      `json:"vResultCode"`
	MStatus     string      `json:"mstatus"`
	MErrorMsg   string      `json:"merrMsg"`
	OrderInfos  *OrderInfos `json:"orderInfos"`
}

type ProperTransactionInfo struct {
	CardTransactionType string `json:"cardTransactionType"`
	ReqWithCapture      string `json:"reqWithCapture"`
	ReqJPOInformation   string `json:"reqJpoInformation"`
}

type TransactionInfo struct {
	Amount      string                `json:"amount"`
	Command     string                `json:"command"`
	MStatus     string                `json:"mstatus"`
	ProperInfo  ProperTransactionInfo `json:"properTransactionInfo"`
	TxnDateTime string                `json:"txnDatetime"`
	TxnID       string                `json:"txnId"`
	VResultCode string                `json:"vResultCode"`
}

type TransactionInfos struct {
	TransactionInfo []TransactionInfo `json:"transactionInfo"`
}

type OrderInfo struct {
	AccountID          string            `json:"accountId"`
	Index              int               `json:"index"`
	OrderID            string            `json:"orderId"`
	ServiceTypeCd      string            `json:"serviceTypeCd"`
	LastSuccessTxnType string            `json:"lastSuccessTxnType"`
	TransactionInfos   *TransactionInfos `json:"transactionInfos"`
}

type OrderInfos struct {
	OrderInfo []OrderInfo `json:"orderInfo"`
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

type ConnectionResponse struct {
	PayNowIDResponse *PayNowIDResponse `json:"payNowIdResponse,omitempty"`
	Result           Result            `json:"result"`
}
