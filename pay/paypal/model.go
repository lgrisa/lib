package paypal

import "encoding/json"

type verifyReq struct {
	AuthAlgo         string          `json:"auth_algo"`
	CertUrl          string          `json:"cert_url"`
	TransmissionId   string          `json:"transmission_id"`
	TransmissionSig  string          `json:"transmission_sig"`
	TransmissionTime string          `json:"transmission_time"`
	WebhookId        string          `json:"webhook_id"`
	WebhookEvent     json.RawMessage `json:"webhook_event"`
}

type verifyResp struct {
	VerificationStatus string `json:"verification_status"`
}

type WebhookNotifyResponse struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	Links            []Link           `json:"links"`
	Status           string           `json:"status"`
	Error            string           `json:"error"`
	ErrorDescription string           `json:"error_description"`
	Scope            string           `json:"scope"`
	AccessToken      string           `json:"access_token"`
	TokenType        string           `json:"token_type"`
	AppId            string           `json:"app_id"`
	ExpiresIn        string           `json:"expires_in"`
	Nonce            string           `json:"nonce"`
	EventType        string           `json:"event_type"`
	Resource         ResponseResource `json:"resource"`
	PurchaseUnits    []PurchaseUnit   `json:"purchase_units"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type ResponseResource struct {
	Id                string                            `json:"id"`
	Status            string                            `json:"status"`
	SupplementaryData ResponseResourceSupplementaryData `json:"supplementary_data"`
}

type ResponseResourceSupplementaryData struct {
	RelatedIds map[string]string `json:"related_ids"`
}

type PurchaseUnit struct {
	Payments map[string][]Capture `json:"payments"`
}

type Capture struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
