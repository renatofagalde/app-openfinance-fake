package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
)

const consentimentoTable = "openfinance-fake-consentimento"

type ConsentimentoRepository struct {
	client *dynamodb.Client
}

func NewConsentimentoRepository(client *dynamodb.Client) *ConsentimentoRepository {
	return &ConsentimentoRepository{client: client}
}

func (r *ConsentimentoRepository) FindByConsentId(consentId string) (*domain.Consentimento, error) {
	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(consentimentoTable),
		Key: map[string]types.AttributeValue{
			"consent_id": &types.AttributeValueMemberS{Value: consentId},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var c domain.Consentimento
	if err := attributevalue.UnmarshalMap(out.Item, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ConsentimentoRepository) Save(c domain.Consentimento) error {
	item, err := attributevalue.MarshalMap(c)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(consentimentoTable),
		Item:      item,
	})
	return err
}

func (r *ConsentimentoRepository) FindAll() ([]domain.Consentimento, error) {
	out, err := r.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(consentimentoTable),
	})
	if err != nil {
		return nil, err
	}

	var items []domain.Consentimento
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ConsentimentoRepository) UpdateStatus(consentId string, status string) error {
	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(consentimentoTable),
		Key: map[string]types.AttributeValue{
			"consent_id": &types.AttributeValueMemberS{Value: consentId},
		},
		UpdateExpression: aws.String("SET consent_status = :s"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: status},
		},
	})
	return err
}
