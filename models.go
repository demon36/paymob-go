package paymob

type Item struct {
	Name        string `json:"name"`
	AmountCents uint   `json:"amount_cents"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
}

type BillingData struct {
	Apartment      string `json:"apartment"`
	Email          string `json:"email"`
	Floor          string `json:"floor"`
	FirstName      string `json:"first_name"`
	Street         string `json:"street"`
	Building       string `json:"building"`
	PhoneNumber    string `json:"phone_number"`
	ShippingMethod string `json:"shipping_method"`
	PostalCode     string `json:"postal_code"`
	City           string `json:"city"`
	Country        string `json:"country"`
	LastName       string `json:"last_name"`
	State          string `json:"state"`
}

type OrderRegistrationRequest struct {
	AuthToken      string `json:"auth_token"`
	DeliveryNeeded bool   `json:"delivery_needed"`
	AmountCents    uint   `json:"amount_cents"`
	Items          []Item `json:"items"`
}

type PaymentKeyRequest struct {
	AuthToken       string      `json:"auth_token"`
	AmountCents     uint        `json:"amount_cents"`
	ExpirationMS    uint        `json:"expiration"`
	OrderID         int         `json:"order_id"`
	Currency        string      `json:"currency"`
	IntegrationID   string      `json:"integration_id"`
	UserBillingData BillingData `json:"billing_data"`
}

//TODO: find a way to merge these two structs
//for post-payment redirection
type TransactionResponseRequest struct {
	Id                   uint   `schema:"id"`
	Pending              bool   `schema:"pending"`
	AmountCents          uint   `schema:"amount_cents"`
	Success              bool   `schema:"success"`
	IsAuth               bool   `schema:"is_auth"`
	IsCapture            bool   `schema:"is_capture"`
	IsStandalonePayment  bool   `schema:"is_standalone_payment"`
	IsVoided             bool   `schema:"is_voided"`
	IsRefunded           bool   `schema:"is_refunded"`
	Is3dSecure           bool   `schema:"is_3d_secure"`
	IntegrationId        uint   `schema:"integration_id"`
	HasParentTransaction bool   `schema:"has_parent_transaction"`
	OrderId              uint   `schema:"order"`
	CreatedAt            string `schema:"created_at"`
	Currency             string `schema:"currency"`
	SourceData struct {
		Pan        string `schema:"pan"`
		Type       string `schema:"type"`
		SubType    string `schema:"sub_type"`
	} `schema:"source_data"`
	ErrorCccured         bool   `schema:"error_occured"`
	Owner                string `schema:"owner"`
}

//for back-end direct callback
type TransactionProcessedRequest struct {
	Obj struct {
		Id                   uint `json:"id"`
		Pending              bool `json:"pending"`
		AmountCents          uint `json:"amount_cents"`
		Success              bool `json:"success"`
		IsAuth               bool `json:"is_auth"`
		IsCapture            bool `json:"is_capture"`
		IsStandalonePayment  bool `json:"is_standalone_payment"`
		IsVoided             bool `json:"is_voided"`
		IsRefunded           bool `json:"is_refunded"`
		Is3dSecure           bool `json:"is_3d_secure"`
		IntegrationId        uint `json:"integration_id"`
		HasParentTransaction bool `json:"has_parent_transaction"`
		Order                struct {
			Id uint `json:"id"`
		} `json:"order"`
		CreatedAt  string `json:"created_at"`
		Currency   string `json:"currency"`
		SourceData struct {
			Pan     string `json:"pan"`
			Type_   string `json:"type"`
			SubType string `json:"sub_type"`
		} `json:"source_data"`
		ErrorOccured bool `json:"error_occured"`
		Owner        uint `json:"owner"`
	} `json:"obj"`
}
