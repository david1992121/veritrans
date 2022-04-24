package veritrans

import (
	"errors"
	"fmt"
)

type PaymentService struct {
	Config ConnectionConfig
}

func NewPaymentService(config ConnectionConfig) *PaymentService {
	return &PaymentService{Config: config}
}

// Get connection paramater from params
func (pay PaymentService) getConnectionParam(param *Params) (*ConnectionParam, error) {
	payNowIDParam := &param.PayNowIDParam
	payNowIDParam.Default()

	connectionParam := &ConnectionParam{
		Params: Params{
			OrderID:       param.OrderID,
			Amount:        param.Amount,
			JPO:           "10",
			WithCapture:   param.WithCapture,
			PayNowIDParam: param.PayNowIDParam,
			TxnVersion:    pay.Config.TxnVersion,
			DummyRequest:  pay.Config.DummyRequest,
			MerchantCCID:  pay.Config.MerchantCCID,
		},
		AuthHash: "",
	}

	if err := SetHash(connectionParam, pay.Config.MerchantCCID, pay.Config.MerchantPassword); err != nil {
		return nil, err
	}
	return connectionParam, nil
}

// Execute Payment
func (pay PaymentService) executePaymentProcess(serviceType PaymentServiceType, mode PaymentManagementMode, param *Params) (bool, error) {
	connectionParam, err := pay.getConnectionParam(param)
	if err == nil {
		paymentRes, err := ProcessRequest(
			fmt.Sprintf("%s/%s/%s", pay.Config.AccountApiURL, PaymentManagementModes[mode], PaymentServiceTypes[serviceType]), connectionParam)
		if err == nil {
			if paymentRes.Result.MStatus == "success" {
				return true, nil
			}

			return false, errors.New(paymentRes.Result.MErrorMsg)
		}
	}
	return false, err
}

// Authorize function
func (pay PaymentService) Authorize(param *Params, serviceType PaymentServiceType) (bool, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodAuthorize),
		param)
}

// Capture function
func (pay PaymentService) Capture(param *Params, serviceType PaymentServiceType) (bool, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodCapture),
		param)
}

// Cancel function
func (pay PaymentService) Cancel(param *Params, serviceType PaymentServiceType) (bool, error) {
	return pay.executePaymentProcess(
		serviceType,
		PaymentManagementMode(MethodCancel),
		param)
}
