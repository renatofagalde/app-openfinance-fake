package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
)

const permissaoTable = "openfinance-fake-permissao"

type PermissaoRepository struct {
	client *dynamodb.Client
}

func NewPermissaoRepository(client *dynamodb.Client) *PermissaoRepository {
	return &PermissaoRepository{client: client}
}

func (r *PermissaoRepository) FindByConsentId(consentId string) ([]domain.Permissao, error) {
	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(permissaoTable),
		KeyConditionExpression: aws.String("consent_id = :cid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cid": &types.AttributeValueMemberS{Value: consentId},
		},
	})
	if err != nil {
		return nil, err
	}

	var items []domain.Permissao
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PermissaoRepository) FindByConsentIdAndPermission(consentId string, permission string) (*domain.Permissao, error) {
	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(permissaoTable),
		Key: map[string]types.AttributeValue{
			"consent_id": &types.AttributeValueMemberS{Value: consentId},
			"permission": &types.AttributeValueMemberS{Value: permission},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var p domain.Permissao
	if err := attributevalue.UnmarshalMap(out.Item, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PermissaoRepository) Save(p domain.Permissao) error {
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(permissaoTable),
		Item:      item,
	})
	return err
}

func (r *PermissaoRepository) UpdateLancar403(consentId string, permission string, lancar403 bool) error {
	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(permissaoTable),
		Key: map[string]types.AttributeValue{
			"consent_id": &types.AttributeValueMemberS{Value: consentId},
			"permission": &types.AttributeValueMemberS{Value: permission},
		},
		UpdateExpression: aws.String("SET lancar_403 = :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberBOOL{Value: lancar403},
		},
	})
	return err
}
