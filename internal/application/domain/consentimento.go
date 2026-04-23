package domain

type Consentimento struct {
	ConsentId     string `dynamodbav:"consent_id"`
	ConsentStatus string `dynamodbav:"consent_status"`
}

const (
	StatusAuthorised = "AUTHORISED"
	StatusRejected   = "REJECTED"
)
