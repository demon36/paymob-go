package paymob

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

	return res["id"].(int), nil
}

func RequestPaymentKey(authToken string, paymentIntegrationId string, orderId int, amountCents uint, firstName string, lastName string, email string, phone string) (string, error) {
	paymentKeyReqData := PaymentKeyRequest{
		AuthToken:    authToken,
		AmountCents:  amountCents,
		ExpirationMS: 600000,
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

func ConcatTransactionResponseValues(v *url.Values) string {
	return v.Get("amount_cents") +
		v.Get("created_at") +
		v.Get("currency") +
		v.Get("error_occured") +
		v.Get("has_parent_transaction") +
		v.Get("id") +
		v.Get("integration_id") +
		v.Get("is_3d_secure") +
		v.Get("is_auth") +
		v.Get("is_capture") +
		v.Get("is_refunded") +
		v.Get("is_standalone_payment") +
		v.Get("is_voided") +
		v.Get("order") +
		v.Get("owner") +
		v.Get("pending") +
		v.Get("source_data.pan") +
		v.Get("source_data.sub_type") +
		v.Get("source_data.type") +
		v.Get("success")
}

func uitoa(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func ConcatTransactionProcessedValues(c *TransactionProcessedRequest) string {
	neededValues :=
		uitoa(c.Obj.Amount_cents) +
			c.Obj.Created_at +
			c.Obj.Currency +
			strconv.FormatBool(c.Obj.Error_occured) +
			strconv.FormatBool(c.Obj.Has_parent_transaction) +
			uitoa(c.Obj.Id) +
			uitoa(c.Obj.Integration_id) +
			strconv.FormatBool(c.Obj.Is_3d_secure) +
			strconv.FormatBool(c.Obj.Is_auth) +
			strconv.FormatBool(c.Obj.Is_capture) +
			strconv.FormatBool(c.Obj.Is_refunded) +
			strconv.FormatBool(c.Obj.Is_standalone_payment) +
			strconv.FormatBool(c.Obj.Is_voided) +
			uitoa(c.Obj.Order.Id) +
			uitoa(c.Obj.Owner) +
			strconv.FormatBool(c.Obj.Pending) +
			c.Obj.Source_data.Pan +
			c.Obj.Source_data.Sub_type +
			c.Obj.Source_data.Type_ +
			strconv.FormatBool(c.Obj.Success)
	return neededValues
}

func ValidateHMAC(message string, expectedHMAC string, key string) bool {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(message))
	calculatedHMAC := mac.Sum(nil)
	expectedHMACBytes, _ := hex.DecodeString(expectedHMAC)
	log.Printf("calculatedHMAC: %x", calculatedHMAC)
	log.Printf("expectedHMAC: %x", expectedHMACBytes)
	return hmac.Equal(expectedHMACBytes, calculatedHMAC)
}
