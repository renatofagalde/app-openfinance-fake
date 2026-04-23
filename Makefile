APP_NAME=openfinance-fake
DIST_DIR=dist
BINARY=$(DIST_DIR)/bootstrap
TEMPLATE=deployments/sam/template.yaml
PROFILE=api-dev
REGION=us-east-1

.PHONY: all clean build deploy stack logs

# Build + deploy do código (uso cotidiano)
all: clean build deploy

clean:
	rm -rf $(DIST_DIR) deploy.zip

build:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) cmd/lambda/main.go
	chmod +x $(BINARY)

deploy:
	zip -j deploy.zip $(BINARY)
	aws lambda update-function-code \
		--function-name $(APP_NAME) \
		--zip-file fileb://deploy.zip \
		--region $(REGION) \
		--profile $(PROFILE)
	rm -f deploy.zip

# Apenas na primeira vez — cria a stack completa (Lambda + API Gateway + DynamoDB)
stack:
	sam build --template $(TEMPLATE)
	sam deploy \
		--template $(TEMPLATE) \
		--stack-name $(APP_NAME) \
		--capabilities CAPABILITY_IAM \
		--region $(REGION) \
		--profile $(PROFILE) \
		--resolve-s3 \
		--no-confirm-changeset

logs:
	aws logs tail /aws/lambda/$(APP_NAME) --follow --profile $(PROFILE)