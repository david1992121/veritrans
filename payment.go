package veritrans

import (
	"errors"
	"fmt"
)

// PaymentService is a service for the payment api
type PaymentService struct {
	Config ConnectionConfig
}

// NewPaymentService initializes the PaymentService
func NewPaymentService(config ConnectionConfig) (*PaymentService, error) {
	if config.PaymentAPIURL != "" {
		return &PaymentService{Config: config}, nil
	}
	return nil, errors.New("api URL not provided")
}

// Get connection paramater from params
func (pay PaymentService) getConnectionParam(param *Params) (*ConnectionParam, error) {
	newParam := *param
	newParam.TxnVersion = pay.Config.TxnVersion
	newParam.DummyRequest = pay.Config.DummyRequest
	newParam.MerchantCCID = pay.Config.MerchantCCID

	connectionParam := &ConnectionParam{
		Params:   newParam,
		AuthHash: "",
	}

	if err := SetHash(connectionParam, pay.Config.MerchantCCID, pay.Config.MerchantPassword); err != nil {
		return nil, err
	}
	return connectionParam, nil
}

// Execute Payment
func (pay PaymentService) executePaymentProcess(serviceType PaymentServiceType, mode PaymentManagementMode, param *Params) (*Result, error) {
	connectionParam, err := pay.getConnectionParam(param)
	if err != nil {
		return nil, err
	}

	apiURL := pay.Config.PaymentAPIURL
	if mode == PaymentManagementMode(MethodSearch) {
		apiURL = pay.Config.SearchAPIURL
	}
	paymentRes, err := ProcessRequest(
		fmt.Sprintf("%s/%s/%s", apiURL, PaymentManagementModes[mode], PaymentServiceTypes[serviceType]), connectionParam)
	if err != nil {
		return nil, err
	}

	if paymentRes.Result.MStatus == "success" {
		return &paymentRes.Result, nil
	}
	return nil, errors.New(paymentRes.Result.MErrorMsg)
}

// Authorize function
func (pay PaymentService) Authorize(param *Params, serviceType PaymentServiceType) (*Result, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodAuthorize),
		param)
}

// Capture function
func (pay PaymentService) Capture(param *Params, serviceType PaymentServiceType) (*Result, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodCapture),
		param)
}

// Cancel function
func (pay PaymentService) Cancel(param *Params, serviceType PaymentServiceType) (*Result, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodCancel),
		param)
}

// Search function
func (pay PaymentService) Search(param *Params, serviceType PaymentServiceType) (*Result, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodSearch),
		param,
	)
}
