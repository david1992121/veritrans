package veritrans

import (
	"errors"
	"fmt"
)

type AccountService struct {
	Config ConnectionConfig
}

func NewAccountService(config ConnectionConfig) *AccountService {
	return &AccountService{Config: config}
}

// Get connection paramater from account parameter
func (acc AccountService) getConnectionParam(accountParam *AccountParam) (*ConnectionParam, error) {
	payNowIdParam := &PayNowIDParam{
		AccountParam: accountParam,
	}
	payNowIdParam.Default()

	connectionParam := &ConnectionParam{
		Params: Params{
			PayNowIDParam: payNowIdParam,
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
		fmt.Sprintf("%s/%s/%s", acc.Config.AccountApiURL, AccountManagementModes[mode], AccountServiceTypes[serviceType]), connectionParam)
	if err != nil {
		return nil, err
	}

	if accountRes.Result.MStatus == "success" {
		return &accountRes.PayNowIDResponse.Account, nil
	}

	return nil, errors.New(accountRes.Result.MErrorMsg)
}

// Create a veritrans account
func (acc AccountService) CreateAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodAdd),
		accountParam)
}

// Remove a veritrans account
func (acc AccountService) DeleteAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodDelete),
		accountParam)
}

// Get a veritrans account
func (acc AccountService) GetAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodGet),
		accountParam)
}

// Get a veritrans account
func (acc AccountService) RestoreAccount(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(AccountType),
		AccountManagementMode(MethodRestore),
		accountParam)
}

// Create a card
func (acc AccountService) CreateCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodAdd),
		accountParam)
}

// Remove a card
func (acc AccountService) DeleteCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodDelete),
		accountParam)
}

// Update a card information
func (acc AccountService) UpdateCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodUpdate),
		accountParam)
}

// Get a veritrans account
func (acc AccountService) GetCard(accountParam *AccountParam) (*Account, error) {
	return acc.executeAccountProcess(
		AccountServiceType(CardType),
		AccountManagementMode(MethodGet),
		accountParam)
}
