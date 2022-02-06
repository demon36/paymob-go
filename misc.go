package paymob

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
)

//returns a token valid for 1 hour
func Authenticate(APIKey string) (string, error) {
	tokenReqData := map[string]string{"api_key": APIKey}
	jsonData, _ := json.Marshal(tokenReqData)

	resp, err := http.Post(
		"https://accept.paymobsolutions.com/api/auth/tokens",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return res["token"].(string), nil
	} else {
		return "", fmt.Errorf("weaccept returned status %v, resp %v", resp.Status, res)
	}

}

//return thirdparty order identifier
func RegisterOrder(authToken string, items []Item, totalPriceInCents uint) (int, error) {
	orderRegistrationReq := OrderRegistrationRequest{
		AuthToken:      authToken,
		DeliveryNeeded: false,
		AmountCents:    totalPriceInCents,
		Items:          items,
	}
	jsonData, _ := json.Marshal(orderRegistrationReq)
	resp, _ := http.Post(
		"https://accept.paymobsolutions.com/api/ecommerce/orders",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	//TODO: define response model
	var res map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("weaccept returned status %v, resp %v", resp.Status, res)
	}

	return int(res["id"].(float64)), nil
}

func RequestPaymentKey(authToken string, paymentIntegrationId string, orderId int, amountCents uint, firstName string, lastName string, email string, phone string, expirationSec uint) (string, error) {
	paymentKeyReqData := PaymentKeyRequest{
		AuthToken:    authToken,
		AmountCents:  amountCents,
		ExpirationMS: expirationSec,
		OrderID:      orderId,
		UserBillingData: BillingData{
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			PhoneNumber:    phone,
			Apartment:      "NA",
			Floor:          "NA",
			Street:         "NA",
			Building:       "NA",
			ShippingMethod: "NA",
			PostalCode:     "NA",
			City:           "NA",
			Country:        "NA",
			State:          "NA",
		},
		Currency:      "EGP",
		IntegrationID: paymentIntegrationId,
	}

	jsonData, _ := json.Marshal(paymentKeyReqData)
	resp, _ := http.Post(
		"https://accept.paymobsolutions.com/api/acceptance/payment_keys",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	//TODO: define response model
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("weaccept returned status %v, resp %v", resp.Status, res)
	}

	return res["token"].(string), nil
}

func GenerateIFrameURL(iframeId string, paymentKey string) string {
	return fmt.Sprintf("https://accept.paymobsolutions.com/api/acceptance/iframes/%s?payment_token=%s",
		iframeId, paymentKey)
}

var decoder = schema.NewDecoder()

func MakeTransactionResponseRequest(params *map[string][]string) TransactionResponseRequest {
	var r TransactionResponseRequest
	decoder.Decode(&r, *params)
	return r
}

func ConcatTransactionResponseValues(r *TransactionResponseRequest) string {
	return uitoa(r.AmountCents) +
		r.CreatedAt +
		r.Currency +
		strconv.FormatBool(r.ErrorCccured) +
		strconv.FormatBool(r.HasParentTransaction) +
		uitoa(r.Id) +
		uitoa(r.IntegrationId) +
		strconv.FormatBool(r.Is3dSecure) +
		strconv.FormatBool(r.IsAuth) +
		strconv.FormatBool(r.IsCapture) +
		strconv.FormatBool(r.IsRefunded) +
		strconv.FormatBool(r.IsStandalonePayment) +
		strconv.FormatBool(r.IsVoided) +
		uitoa(r.OrderId) +
		r.Owner +
		strconv.FormatBool(r.Pending) +
		r.SourceDataPan +
		r.SourceDataSubType +
		r.SourceDataType +
		strconv.FormatBool(r.Success)
}

func uitoa(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func ConcatTransactionProcessedValues(c *TransactionProcessedRequest) string {
	neededValues :=
		uitoa(c.Obj.AmountCents) +
			c.Obj.CreatedAt +
			c.Obj.Currency +
			strconv.FormatBool(c.Obj.ErrorOccured) +
			strconv.FormatBool(c.Obj.HasParentTransaction) +
			uitoa(c.Obj.Id) +
			uitoa(c.Obj.IntegrationId) +
			strconv.FormatBool(c.Obj.Is3dSecure) +
			strconv.FormatBool(c.Obj.IsAuth) +
			strconv.FormatBool(c.Obj.IsCapture) +
			strconv.FormatBool(c.Obj.IsRefunded) +
			strconv.FormatBool(c.Obj.IsStandalonePayment) +
			strconv.FormatBool(c.Obj.IsVoided) +
			uitoa(c.Obj.Order.Id) +
			uitoa(c.Obj.Owner) +
			strconv.FormatBool(c.Obj.Pending) +
			c.Obj.SourceData.Pan +
			c.Obj.SourceData.SubType +
			c.Obj.SourceData.Type_ +
			strconv.FormatBool(c.Obj.Success)
	return neededValues
}

func ValidateHMAC(message string, expectedHMAC string, key string) bool {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(message))
	calculatedHMAC := mac.Sum(nil)
	expectedHMACBytes, _ := hex.DecodeString(expectedHMAC)
	return hmac.Equal(expectedHMACBytes, calculatedHMAC)
}
