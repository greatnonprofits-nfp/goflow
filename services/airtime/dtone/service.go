package dtone

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/greatnonprofits-nfp/goflow/flows"
	"github.com/greatnonprofits-nfp/goflow/utils/httpx"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type service struct {
	login    string
	apiToken string
	currency string
}

// NewService creates a new DTOne airtime service
func NewService(login, apiToken, currency string) flows.AirtimeService {
	return &service{
		login:    login,
		apiToken: apiToken,
		currency: currency,
	}
}

func (s *service) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	transfer := &flows.AirtimeTransfer{
		Sender:        sender,
		Recipient:     recipient,
		DesiredAmount: decimal.Zero,
		ActualAmount:  decimal.Zero,
	}

	client := NewClient(session.Engine().HTTPClient(), s.login, s.apiToken)

	info, trace, err := client.MSISDNInfo(recipient.Path(), s.currency, "1")
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus))
	}
	if err != nil {
		return transfer, err
	}

	transfer.Currency = info.DestinationCurrency

	// look up the amount to send in this currency
	amount, hasAmount := amounts[transfer.Currency]
	if !hasAmount {
		return transfer, errors.Errorf("no amount configured for transfers in %s", transfer.Currency)
	}

	transfer.DesiredAmount = amount

	// find the product closest to our desired amount
	var useProduct string
	useAmount := decimal.Zero
	for p, product := range info.ProductList {
		price := info.LocalInfoValueList[p]
		if price.GreaterThan(useAmount) && price.LessThanOrEqual(amount) {
			useProduct = product
			useAmount = price
		}
	}

	if useAmount == decimal.Zero {
		return transfer, errors.Errorf("amount requested is smaller than the mimimum topup of %s %s", info.LocalInfoValueList[0].String(), info.DestinationCurrency)
	}

	reservedID, trace, err := client.ReserveID()
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus))
	}
	if err != nil {
		return transfer, err
	}

	topup, trace, err := client.Topup(reservedID.ReservedID, sender.Path(), recipient.Path(), useProduct, "")
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus))
	}
	if err != nil {
		return transfer, err
	}

	transfer.ActualAmount = topup.ActualProductSent

	return transfer, nil
}

func httpLogStatus(t *httpx.Trace) flows.CallStatus {
	// DTOne error responses use HTTP 200 OK but we consider them errors
	if t.ResponseBody != nil {
		base := &baseResponse{}
		unmarshalResponse(t.ResponseBody, base)
		if base.Error() != nil {
			return flows.CallStatusResponseError
		}
	}

	return flows.HTTPStatusFromCode(t)
}
