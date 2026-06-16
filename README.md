# Torque Backend

Backend em Go para gestao de veiculos, dispositivos embarcados e telemetria.

## Visao geral

O projeto tem dois caminhos principais de entrada:

- HTTP API: usada por webapp, backoffice e integracoes sincronas.
- Ingestao assicrona: usada para telemetria e DTCs via fila.

No estado atual deste repositorio:

- a API HTTP vive em `cmd/api`
- o worker de ingestao vive em `cmd/worker`
- a CLI de apoio vive em `cmd/cli`
- a documentacao OpenAPI manual vive em `docs/openapi.yaml`

## Arquitetura de pastas

### `cmd/`

Entrypoints executaveis do sistema.

- `cmd/api`: servidor HTTP principal.
- `cmd/worker`: consumo de mensagens RabbitMQ e escrita em TimescaleDB.
- `cmd/cli`: comandos para publicar eventos de teste e simular sessao de telemetria.

### `internal/core/`

Blocos compartilhados e transversais.

- `appctx`: contexto de autenticacao da request.
- `apperr`: padrao de erro da aplicacao.
- `db`: conexao e tipos compartilhados de banco.
- `id`: geracao de IDs.
- `logger`: setup de logging.
- `middleware`: middleware generico de logging.
- `pagination`: helpers de paginacao.
- `pki`: cliente Step CA.

### `internal/infrastructure/`

Adaptadores entre modulos e integracoes externas.

- `adapters`: resolve dependencias cruzadas entre modulos sem acoplamento direto.
- `messaging`: contratos das mensagens publicadas/consumidas na fila.

### `internal/modules/`

Modulos de negocio. Cada modulo segue, em geral, esta divisao:

- `domain`: entidades, contratos e regras centrais.
- `application/dto`: formatos de entrada e saida.
- `application/usecase`: casos de uso.
- `infrastructure/repository`: implementacoes concretas de persistencia.

Modulos atuais:

- `vehicle`: cadastro de veiculos e catalogo de modelos/anos/cores.
- `device`: cadastro de dispositivos, comissionamento e certificados.
- `telemetry`: leitura e escrita de telemetria e DTCs.
- `dealership`: estrutura reservada; hoje nao tem implementacao ativa.

### `migrations/`

Migracoes SQL separadas por banco:

- `migrations/main`: schema transacional no Postgres.
- `migrations/timescale`: schema de series temporais no TimescaleDB.

### `seeds/`

Dados de referencia para ambiente local:

- catalogo de modelos
- veiculos de teste
- catalogo de DTCs

### `infra/`

Infraestrutura local via Docker Compose.

Servicos ativos no modo dev atual:

- `postgres`
- `timescaledb`
- `rabbitmq`
- `step-ca`
- `minio`
- `mailhog`
- `traefik`
- `oathkeeper`
- `kratos`
- `kratos-postgres`

### `docs/`

Contrato OpenAPI e suporte a documentacao.

- `docs/openapi.yaml`: especificacao manual da API.
- `docs/docs.go`: embed da spec no binario.

## Modulos funcionais

### Veiculos

Responsavel por:

- criar, listar, buscar, atualizar e remover veiculos
- listar modelos de veiculo
- listar anos e cores por modelo
- enriquecer a resposta do veiculo com status de DTC ativo

### Dispositivos

Responsavel por:

- criar dispositivos
- emitir certificados via Step CA
- listar dispositivos
- comissionar dispositivo em veiculo
- descomissionar dispositivo

### Telemetria

Responsavel por:

- gravar telemetria e DTCs recebidos assincronamente
- listar telemetria por veiculo
- listar DTCs ativos
- consultar TimescaleDB

## Fluxo de execucao

### API HTTP

Fluxo principal:

`webapp/client -> API -> modules -> Postgres/Timescale/Step CA`

Rotas de documentacao:

- `GET /docs`
- `GET /openapi.yaml`

### Worker

Fluxo principal:

`RabbitMQ -> worker -> telemetry module -> TimescaleDB`

## Autenticacao em desenvolvimento

O middleware de auth espera:

- `x-user-id`
- `x-user-role`

Para facilitar desenvolvimento local, existe um bypass por `.env` em `cmd/api/.env`.

Variaveis:

```env
APP_ENV=development
AUTH_BYPASS_ENABLED=false
AUTH_BYPASS_USER_ID=11111111-1111-1111-1111-111111111111
AUTH_BYPASS_USER_ROLE=admin
```

Regras:

- se `x-user-id` vier no header, ele tem prioridade
- se o header nao vier e `AUTH_BYPASS_ENABLED=true`, a API usa os valores do `.env`
- o bypass so e aceito quando `APP_ENV=development`

## Rodando localmente

Suba a infraestrutura e rode migrations/seeds:

```bash
./scripts/setup.sh
```

Isso sobe a stack local com `docker compose`, aplica migrations, prepara o Step CA e carrega dados de referencia.

Depois, rode a API e o worker separadamente, se necessario:

```bash
go run ./cmd/api
go run ./cmd/worker
```

## Observacoes

- `mocks/` contem mocks usados pelos testes unitarios.
- `docs/openapi.yaml` deve ser mantido junto com as mudancas de contrato da API.
- o compose de dev sobe Traefik, Oathkeeper e Kratos para o fluxo local de autenticacao.
