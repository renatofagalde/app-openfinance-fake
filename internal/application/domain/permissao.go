package domain

type Permissao struct {
	ConsentId  string `dynamodbav:"consent_id"`
	Permission string `dynamodbav:"permission"`
	Lancar403  bool   `dynamodbav:"lancar_403"`
}
