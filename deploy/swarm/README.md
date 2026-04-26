# Docker Swarm com Secrets (stacks separadas)

Estrutura para deploy em duas stacks independentes:

- `campus-db` com PostgreSQL (`db-stack.yml`)
- `campus-api` com a API (`api-stack.yml`)

As credenciais ficam em `docker secrets`.

## 1) Inicializar swarm (se necessario)

```bash
docker swarm init
```

## 2) Build da imagem da API

Na raiz do projeto:

```bash
docker build -t campus_connect_api:latest .
```

Se for usar em outro node/cluster, publique em registry e ajuste `API_IMAGE`.

## 3) Criar os secrets

```bash
printf "SENHA_FORTE_AQUI" | docker secret create postgres_password -
printf "postgres://campus_connect:SENHA_FORTE_AQUI@db:5432/campus_connect?sslmode=disable" | docker secret create database_url -
```

## 4) Preparar variaveis

```bash
cp .env.db.example .env.db
cp .env.api.example .env.api
```

## 5) Deploy das stacks

Primeiro banco (cria a rede overlay compartilhada `campus_overlay_shared`):

```bash
docker stack deploy -c db-stack.yml --env-file .env.db campus-db
```

Depois API:

```bash
docker stack deploy -c api-stack.yml --env-file .env.api campus-api
```

## 6) Verificacao

```bash
docker stack services campus-db
docker stack services campus-api
docker service logs -f campus-api_api
```

## Observacoes

- O host do banco dentro da rede overlay e `db`.
- A API usa `DATABASE_URL_FILE=/run/secrets/database_url`.
- Secret no Swarm nao e editavel; remova e recrie quando precisar trocar valor.
