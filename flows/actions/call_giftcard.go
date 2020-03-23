package actions

import (
	"encoding/json"
	"github.com/greatnonprofits-nfp/goflow/flows"
	"github.com/greatnonprofits-nfp/goflow/flows/events"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func init() {
	registerType(TypeCallGiftcard, func() flows.Action { return &CallGiftcardAction{} })
}

// TypeCallGiftcard is the type for the call giftcard action
const TypeCallGiftcard string = "call_giftcard"

// CallGiftcardAction can be used to call an external service. The body, header and url fields may be
// templates and will be evaluated at runtime. A [event:giftcard_called] event will be created based on
// the results of the HTTP call. If this action has a `result_name`, then addtionally it will create
// a new result with that name. If the lookup returned valid JSON, that will be accessible
// through `extra` on the result.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "call_giftcard",
//     "giftcard_db": {"id": "demo_test_giftcard", "text": "Test Giftcard"},
//     "result_name": "giftcard"
//   }
//
// @action call_giftcard
type CallGiftcardAction struct {
	baseAction
	onlineAction

	DB           map[string]string `json:"giftcard_db"`
	GiftcardType string            `json:"giftcard_type"`
	ResultName   string            `json:"result_name,omitempty"`
}

// NewCallGiftcardAction creates a new call giftcard action
func NewCallGiftcardAction(uuid flows.ActionUUID, giftcardDb map[string]string, giftcardType string, resultName string) *CallGiftcardAction {
	return &CallGiftcardAction{
		baseAction:   newBaseAction(TypeCallGiftcard, uuid),
		DB:           giftcardDb,
		GiftcardType: giftcardType,
		ResultName:   resultName,
	}
}

// Validate validates our action is valid
func (a *CallGiftcardAction) Validate() error {
	if a.DB["id"] == "" {
		return errors.Errorf("id is required on Giftcard DB")
	}

	return nil
}

// Execute runs this action
func (a *CallGiftcardAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	method := "POST"

	// substitute any variables in our url
	parseUrl := getEnv(envVarServerUrl, "http://localhost:9090/parse")
	var giftcardType string
	if a.GiftcardType == giftcardCheckType {
		giftcardType = "giftcards_remaining"
	} else {
		giftcardType = "giftcard"
	}
	url := parseUrl + "/functions/" + giftcardType

	if parseUrl == "" {
		logEvent(events.NewErrorf("Parse Server URL is an empty string, skipping"))
		return nil
	}

	if a.DB["id"] == "" {
		logEvent(events.NewErrorf("Parse Server DB is required, skipping"))
		return nil
	}

	contact_urn := run.Contact().PreferredURN()

	body := make(map[string]interface{})
	body["db"] = a.DB["id"]
	body["urn"] = contact_urn.URN().Path()

	fdlDefaultURL := getEnv(envVarFDLDefaultURL, "")
	fdlKey := getEnv(envVarFDLKey, "")

	body["fdl_default_url"] = fdlDefaultURL
	body["fdl_key"] = fdlKey

	b, _ := json.Marshal(body)

	return a.call(run, step, url, method, string(b), logEvent)
}

// Execute runs this action
func (a *CallGiftcardAction) call(run flows.FlowRun, step flows.Step, url, method, body string, logEvent flows.EventCallback) error {
	// build our request
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err
	}

	appId := getEnv(envVarAppId, "myAppId")
	masterKey := getEnv(envVarMasterKey, "myMasterKey")

	// add the custom headers, substituting any template vars
	req.Header.Add(xParseApplicationId, appId)
	req.Header.Add(xParseMasterKey, masterKey)
	req.Header.Add("Content-Type", "application/json")

	svc, err := run.Session().Engine().Services().Webhook(run.Session())
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	call, err := svc.Call(run.Session(), req)

	if err != nil {
		logEvent(events.NewError(err))
	}
	if call != nil {
		a.updateWebhook(run, call)

		status := callStatus(call, err, false)

		logEvent(events.NewGiftcardCalled(call, status, ""))

		if a.ResultName != "" {
			a.saveWebhookResult(run, step, a.ResultName, call, status, logEvent)
		}
	}

	return nil
}

// Results enumerates any results generated by this flow object
func (a *CallGiftcardAction) Results(include func(*flows.ResultInfo)) {
	if a.ResultName != "" {
		include(flows.NewResultInfo(a.ResultName, webhookCategories))
	}
}
