package veritrans

// EnvVariables is a list of the environment variables
var EnvVariables = []string{
	"MDK_API_TOKEN",
	"MDK_API_URL",
	"MERCHANT_CCID",
	"MERCHANT_PASSWORD",
	"DUMMY_REQUEST",
	"TXN_VERSION",
	"ACCOUNT_API_URL",
	"PAYMENT_API_URL",
	"SEARCH_API_URL",
}

// ConnectionConfig is a configuration of veritrans connection
// AccountAPIURL is the account management api endpoint (https://api.veritrans.co.jp:443/paynowid/v1/)
// PaymentAPIURL is the payment api endpoint (https://api.veritrans.co.jp:443/paynow/v2)
// TxnVersion is the version of the veritrans api (2.0.0)
// DummyRequest is the flag indicating whether the request is dummy or live
type ConnectionConfig struct {
	MerchantCCID     string
	MerchantPassword string
	AccountAPIURL    string
	PaymentAPIURL    string
	SearchAPIURL     string
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

// RecurringChargeParam represents the "recurringChargeParam" of the request.
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

// PayNowIDParam represents the "payNowIDParam" of the request.
type PayNowIDParam struct {
	Token        string        `json:"token,omitempty"`
	AccountParam *AccountParam `json:"accountParam,omitempty"`
	Memo         string        `json:"memo1,omitempty"`
	FreeKey      string        `json:"freeKey,omitempty"`
}

// OrderParam struct
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
	PayNowIDParam    *PayNowIDParam `json:"payNowIDParam,omitempty"`
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

// AccountManagementMode is the enum type for the account management mode
type AccountManagementMode int32

const (
	// MethodAdd indicates a Add method
	MethodAdd AccountManagementMode = iota
	// MethodUpdate indicates a Update method
	MethodUpdate
	// MethodDelete indicates a Delete method
	MethodDelete
	// MethodRestore indicates a Restore method
	MethodRestore
	// MethodGet indicates a Get method
	MethodGet
)

// AccountManagementModes list
var AccountManagementModes = []string{"Add", "Update", "Delete", "Restore", "Get"}

// AccountServiceType is the enum type of provided services
type AccountServiceType int32

const (
	// AccountType represents the account service
	AccountType AccountServiceType = iota
	// CardType represents the card service
	CardType
)

// AccountServiceTypes is list of services
var AccountServiceTypes = []string{"account", "cardinfo"}

// PaymentManagementMode is the enum type for the payment mode
type PaymentManagementMode int32

const (
	// MethodAuthorize indicates a Authorize
	MethodAuthorize PaymentManagementMode = iota
	// MethodCapture indicates a Capture
	MethodCapture
	// MethodCancel indicates a Cancel
	MethodCancel
	// MethodSearch indicates a Search
	MethodSearch
)

// PaymentManagementModes a list of methods
var PaymentManagementModes = []string{"Authorize", "Capture", "Cancel", "Search"}

// PaymentServiceType represents the payment service
type PaymentServiceType int32

const (
	// PayCard indicates the "card"
	PayCard PaymentServiceType = iota
	// MPI indicates the "mpi"
	MPI
	// CVS indicates the "cvs"
	CVS
	// EM indicates the "em"
	EM
	// Bank indicates the "bank"
	Bank
	// UPop indicates the "upop"
	UPop
	// Paypal indicates the "paypal"
	Paypal
	// Saison indicates the "saison"
	Saison
	// Alipay indicates the "alipay"
	Alipay
	// Carrier indicates the "carrier"
	Carrier
	// Search indicates the "search"
	Search
)

// PaymentServiceTypes is a list of services
var PaymentServiceTypes = []string{"card", "mpi", "cvs", "em", "bank", "upop", "paypal", "saison", "alipay", "carrier", "search"}

// Default function for the PayNowIDParam
func (payParam *PayNowIDParam) Default() {
	if payParam.Memo == "" {
		payParam.Memo = "memo"
	}
	if payParam.FreeKey == "" {
		payParam.FreeKey = "freekey"
	}
}

// Default function for the AccountBasicParam
func (accountBasicParam *AccountBasicParam) Default() {
	if accountBasicParam.ForceDeleteDate == "" {
		accountBasicParam.ForceDeleteDate = "0"
	}
}

// Default function for the RecurringChargeParam
func (recurringChargeParam *RecurringChargeParam) Default() {
	if recurringChargeParam.FinalCharge == "" {
		recurringChargeParam.FinalCharge = "0"
	}
}

// Result indicates the response of api
type Result struct {
	VResultCode string      `json:"vResultCode"`
	MStatus     string      `json:"mstatus"`
	MErrorMsg   string      `json:"merrMsg"`
	OrderInfos  *OrderInfos `json:"orderInfos"`
}

// ProperTransactionInfo struct
type ProperTransactionInfo struct {
	CardTransactionType string `json:"cardTransactionType"`
	ReqWithCapture      string `json:"reqWithCapture"`
	ReqJPOInformation   string `json:"reqJpoInformation"`
}

// TransactionInfo struct
type TransactionInfo struct {
	Amount      string                `json:"amount"`
	Command     string                `json:"command"`
	MStatus     string                `json:"mstatus"`
	ProperInfo  ProperTransactionInfo `json:"properTransactionInfo"`
	TxnDateTime string                `json:"txnDatetime"`
	TxnID       string                `json:"txnId"`
	VResultCode string                `json:"vResultCode"`
}

// TransactionInfos struct
type TransactionInfos struct {
	TransactionInfo []TransactionInfo `json:"transactionInfo"`
}

// OrderInfo struct
type OrderInfo struct {
	AccountID          string            `json:"accountId"`
	Index              int               `json:"index"`
	OrderID            string            `json:"orderId"`
	ServiceTypeCd      string            `json:"serviceTypeCd"`
	LastSuccessTxnType string            `json:"lastSuccessTxnType"`
	TransactionInfos   *TransactionInfos `json:"transactionInfos"`
}

// OrderInfos struct
type OrderInfos struct {
	OrderInfo []OrderInfo `json:"orderInfo"`
}

// CardInfo struct
type CardInfo struct {
	CardExpire  string `json:"cardExpire"`
	CardID      string `json:"cardId"`
	CardNumber  string `json:"cardNumber"`
	DefaultCard string `json:"defaultCard"`
}

// Account struct
type Account struct {
	AccountID string     `json:"accountId"`
	CardInfo  []CardInfo `json:"cardInfo"`
}

// PayNowIDResponse struct
type PayNowIDResponse struct {
	Account Account `json:"account"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

// ConnectionResponse struct
type ConnectionResponse struct {
	PayNowIDResponse *PayNowIDResponse `json:"payNowIdResponse,omitempty"`
	Result           Result            `json:"result"`
}
