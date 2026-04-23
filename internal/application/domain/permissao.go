package domain

type Permissao struct {
	ConsentId  string `dynamodbav:"consent_id"  json:"consent_id"`
	Permission string `dynamodbav:"permission"  json:"permission"`
	Lancar403  bool   `dynamodbav:"lancar_403"  json:"lancar_403"`
}
