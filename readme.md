# app-openfinance-fake

Mock server do Open Finance Brasil criado para fins de estudo e entendimento do ecossistema, a partir da experiência prática trabalhando no setor bancário.

Simula o servidor de autorização e os endpoints de dados de contas, permitindo testar cenários de **sucesso (200)**, **autorização negada (401)** e **permissão negada (403)** de forma controlada.

---

## Como funciona

No Open Finance Brasil, uma instituição receptora consome dados de contas de outras instituições transmissoras. Para isso, precisa:

1. Obter um **token de acesso** via `client_credentials`
2. Verificar se o **consentimento** do usuário está ativo (`AUTHORISED`)
3. Verificar se as **permissões granulares** foram concedidas antes de chamar cada endpoint
4. Tratar corretamente os erros **401** (consentimento inválido) e **403** (permissão negada)

Este projeto simula a transmissora, permitindo testar todo esse fluxo localmente sem depender de ambientes externos.

---

## Arquitetura

```
┌─────────────────────────────────────────────────────┐
│              Robô Receptor (.NET / Go)              │
│         Consome dados via Open Finance              │
└──────────────────────┬──────────────────────────────┘
                       │ TOKEN_URL / API_URL
                       ▼
┌─────────────────────────────────────────────────────┐
│              API Gateway (AWS)                      │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│              Lambda (Go / Gin)                      │
│              app-openfinance-fake                   │
└──────────┬───────────────────────┬──────────────────┘
           │                       │
           ▼                       ▼
┌──────────────────┐   ┌──────────────────────────────┐
│   DynamoDB       │   │   DynamoDB                   │
│   consentimento  │   │   permissao                  │
│   consent_id     │   │   consent_id + permission    │
│   consent_status │   │   lancar_403 (bool)          │
└──────────────────┘   └──────────────────────────────┘
```

---

## Stack

- **Go 1.21** — Gin Router + arquitetura hexagonal
- **AWS Lambda** — `provided.al2023` / `x86_64`
- **AWS API Gateway** — stage `dev`
- **AWS DynamoDB** — PAY_PER_REQUEST
- **AWS SAM** — deploy da infraestrutura

---

## Estrutura do projeto

```
app-openfinance-fake/
├── cmd/lambda/main.go
├── internal/
│   ├── adapter/
│   │   ├── input/
│   │   │   ├── controller/
│   │   │   │   ├── token_controller.go
│   │   │   │   ├── accounts_controller.go
│   │   │   │   └── mock_controller.go
│   │   │   └── routes/routes.go
│   │   └── output/dynamo/
│   │       ├── consentimento_repository.go
│   │       └── permissao_repository.go
│   ├── application/
│   │   ├── domain/
│   │   ├── port/input/
│   │   ├── port/output/
│   │   └── service/
│   └── config/logger/
├── deployments/sam/template.yaml
├── Makefile
└── go.mod
```

---

## Pré-requisitos

- Go 1.21+
- AWS CLI configurado
- AWS SAM CLI (`pip install aws-sam-cli`)

---

## Instalação

```shell
git clone https://github.com/renatofagalde/app-openfinance-fake.git
cd app-openfinance-fake
go mod tidy
```

---

## Deploy

### Primeira vez — cria toda a infraestrutura

```shell
make stack
```

Cria a Lambda, API Gateway e as tabelas DynamoDB. A URL da API é exibida nos outputs ao final.

### Atualizações de código

```shell
make all
```

---

## Endpoints

### Token

| Método | Endpoint |
|--------|----------|
| `POST` | `/auth/realms/{realm}/protocol/openid-connect/token` |

### Open Finance — Contas

| Método | Endpoint | Permission verificada |
|--------|----------|----------------------|
| `GET` | `/accounts` | `ACCOUNTS_READ` |
| `GET` | `/accounts/{accountId}` | `ACCOUNTS_READ` |
| `GET` | `/accounts/{accountId}/balances` | `ACCOUNTS_BALANCES_READ` |
| `GET` | `/accounts/{accountId}/transactions` | `ACCOUNTS_TRANSACTIONS_READ` |
| `GET` | `/accounts/{accountId}/transactions-current` | `ACCOUNTS_TRANSACTIONS_READ` |
| `GET` | `/accounts/{accountId}/overdraft-limits` | `ACCOUNTS_OVERDRAFT_LIMITS_READ` |

### Mock — Manipulação de dados

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/mock/consentimentos` | Lista consentimentos |
| `POST` | `/mock/consentimentos` | Cria consentimento |
| `PUT` | `/mock/consentimentos/{consentId}/status` | Atualiza status |
| `GET` | `/mock/permissoes/{consentId}` | Lista permissions |
| `POST` | `/mock/permissoes` | Cria permission |
| `PUT` | `/mock/permissoes/{consentId}/{permission}/lancar403` | Ativa/desativa 403 |

---

## Cenários de teste

| # | Consent ID | Status | Permission com 403 | Resultado esperado |
|---|------------|--------|-------------------|--------------------|
| C1 | `urn:obc-fake:consent-sucesso-001` | `AUTHORISED` | Nenhuma | Todos os endpoints retornam 200 |
| C2 | `urn:obc-fake:consent-rejeitado-002` | `REJECTED` | N/A | 401 — consentimento inválido |
| C3 | `urn:obc-fake:consent-sem-balances-003` | `AUTHORISED` | `ACCOUNTS_BALANCES_READ` | Balances retorna 403, demais 200 |
| C4 | `urn:obc-fake:consent-sem-transactions-004` | `AUTHORISED` | `ACCOUNTS_TRANSACTIONS_READ` | Transactions retorna 403, demais 200 |
| C5 | `urn:obc-fake:consent-sem-overdraft-005` | `AUTHORISED` | `ACCOUNTS_OVERDRAFT_LIMITS_READ` | Overdraft retorna 403, demais 200 |

---

## Comandos úteis

### Logs

```shell
aws logs tail /aws/lambda/openfinance-fake --follow --profile api-dev
```

### DynamoDB — Listar consentimentos

```shell
aws dynamodb scan \
  --table-name openfinance-fake-consentimento \
  --profile api-dev \
  --query "Items[*].consent_id.S" \
  --output text
```

### DynamoDB — Deletar item específico

```shell
aws dynamodb delete-item \
  --table-name openfinance-fake-consentimento \
  --key '{"consent_id": {"S": "urn:obc-fake:consent-sucesso-001"}}' \
  --profile api-dev
```

### DynamoDB — Listar permissions de um consentimento

```shell
aws dynamodb query \
  --table-name openfinance-fake-permissao \
  --key-condition-expression "consent_id = :cid" \
  --expression-attribute-values '{":cid": {"S": "urn:obc-fake:consent-sucesso-001"}}' \
  --profile api-dev
```

---

## Observações

- As tabelas DynamoDB usam `PAY_PER_REQUEST` — sem custo fixo, ideal para uso esporádico
- O token retornado é sempre estático e não expira
- A URL da API Gateway não está exposta neste repositório — consulte os outputs do CloudFormation após o deploy
- Para zerar os dados entre testes, remova os itens diretamente no DynamoDB via console AWS ou pelos comandos acima