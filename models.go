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
	OrderID         string      `json:"order_id"`
	Currency        string      `json:"currency"`
	IntegrationID   string      `json:"integration_id"`
	UserBillingData BillingData `json:"billing_data"`
}

type TransactionProcessedRequest struct {
	Obj struct {
		Id                     uint `json:"id"`
		Pending                bool `json:"pending"`
		Amount_cents           uint `json:"amount_cents"`
		Success                bool `json:"success"`
		Is_auth                bool `json:"is_auth"`
		Is_capture             bool `json:"is_capture"`
		Is_standalone_payment  bool `json:"is_standalone_payment"`
		Is_voided              bool `json:"is_voided"`
		Is_refunded            bool `json:"is_refunded"`
		Is_3d_secure           bool `json:"is_3d_secure"`
		Integration_id         uint `json:"integration_id"`
		Has_parent_transaction bool `json:"has_parent_transaction"`
		Order                  struct {
			Id uint `json:"id"`
		} `json:"order"`
		Created_at  string `json:"created_at"`
		Currency    string `json:"currency"`
		Source_data struct {
			Pan      string `json:"pan"`
			Type_    string `json:"type"`
			Sub_type string `json:"sub_type"`
		} `json:"source_data"`
		Error_occured bool `json:"error_occured"`
		Owner         uint `json:"owner"`
	} `json:"obj"`
}
