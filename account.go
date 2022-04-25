package veritrans

import (
	"errors"
	"fmt"
)

// AccountService is a service for managing accounts
type AccountService struct {
	Config ConnectionConfig
}

// NewAccountService initializes the account service
func NewAccountService(config ConnectionConfig) *AccountService {
	return &AccountService{Config: config}
}

// Get connection paramater from account parameter
func (acc AccountService) getConnectionParam(accountParam *AccountParam) (*ConnectionParam, error) {
	payNowIDParam := &PayNowIDParam{
		AccountParam: accountParam,
	}
	payNowIDParam.Default()

	connectionParam := &ConnectionParam{
		Params: Params{
			PayNowIDParam: payNowIDParam,
			TxnVersion:    acc.Config.TxnVersion,
			DummyRequest:  acc.Config.DummyRequest,
			MerchantCCID:  acc.Config.MerchantCCID,
		},
		AuthHash: "",
	}

	if err := SetHash(connectionParam, acc.Config.MerchantCCID, acc.Config.MerchantPassword); err != nil {
		return nil, err
	}
	return connectionParam, nil
}

// Execute Account CRUD
func (acc AccountService) executeAccountProcess(serviceType AccountServiceType, mode AccountManagementMode, accountParam *AccountParam) (*Account, error) {
	connectionParam, err := acc.getConnectionParam(accountParam)
	if err != nil {
		return nil, err
	}

	accountRes, err := ProcessRequest(
		fmt.Sprintf("%s/%s/%s", acc.Config.AccountAPIURL, AccountManagementModes[mode], AccountServiceTypes[serviceType]), connectionParam)
	if err != nil {
		return nil, err
	}

	if accountRes.Result.MStatus == "success" {
		return &accountRes.PayNowIDResponse.Account, nil
	}

	return nil, errors.New(accountRes.Result.MErrorMsg)
}

// CreateAccount function
func (acc AccountService) CreateAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodAdd),
		accountParam)
}

// DeleteAccount function
func (acc AccountService) DeleteAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodDelete),
		accountParam)
}

// GetAccount function
func (acc AccountService) GetAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodGet),
		accountParam)
}

// RestoreAccount function
func (acc AccountService) RestoreAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodRestore),
		accountParam)
}

// CreateCard function
func (acc AccountService) CreateCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodAdd),
		accountParam)
}

// DeleteCard function
func (acc AccountService) DeleteCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodDelete),
		accountParam)
}

// UpdateCard function
func (acc AccountService) UpdateCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodUpdate),
		accountParam)
}

// GetCard function
func (acc AccountService) GetCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodGet),
		accountParam)
}
