# Docker Swarm com Secrets (stacks separadas)

Este diretório fica em `infra/swarm`.

Estrutura para deploy em duas stacks independentes:

- `campus-db` com PostgreSQL (`db-stack.yml`)
- `campus-api` com a API (`api-stack.yml`)

As credenciais ficam em `docker secrets`.

## 1) Inicializar swarm (se necessario)

```bash
docker swarm init
```

## 2) Build da imagem da API (otimizado)

Na raiz do projeto, com BuildKit (cache de módulos Go + binário `-s -w`):

```bash
export DOCKER_BUILDKIT=1
docker build -t campus_connect_api:latest .
```

Imagem final: Alpine + binário estático (~25–35 MB comprimida, conforme dependências).
Runtime: `GIN_MODE=release` definido no `api-stack.yml`.

Se for usar em outro node/cluster, publique em registry e ajuste `API_IMAGE`.

### Performance (resumo)

| Componente | Escolha | Motivo |
|------------|---------|--------|
| API build | `CGO_ENABLED=0`, `-ldflags="-s -w"`, `-trimpath` | Binário menor e sem libc |
| API runtime | Alpine + `ca-certificates` + `wget` (health) | Leve; sem shell no distroless impede healthcheck |
| Postgres | `postgres:16-alpine` | Imagem oficial mínima |
| Postgres tuning | `shared_buffers`, `effective_cache_size`, `shm_size` | Ver `.env.db.example` |
| Swarm | `resources.limits` + log rotation | Evita OOM e disco cheio |

**Host 1 GB RAM:** `POSTGRES_SHARED_BUFFERS_MB=64`, `POSTGRES_EFFECTIVE_CACHE_MB=128`, `API_MEM_LIMIT=128M`, `POSTGRES_MEM_LIMIT=384M`.

**Host 4 GB+:** pode subir buffers para 256/512 MB e `API_REPLICAS=2` se tiver load balancer.

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

O `.env.api` e opcional (o YAML tem defaults). Sem ele:

```bash
docker stack deploy -c api-stack.yml campus-api
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
