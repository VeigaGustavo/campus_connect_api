# Docker Swarm com Secrets (stacks separadas)

Este diretório fica em `infra/swarm`.

Estrutura para deploy em duas stacks independentes:

- `campus-db` com PostgreSQL (`db-stack.yml`)
- `campus-api` com a API (`api-stack.yml`)

Credenciais sensíveis ficam em **Docker secrets** (não no `.env`). O `.env` serve só para tuning (RAM, portas, imagem).

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

O `Dockerfile` compila com `-p 1` e `GOMAXPROCS=1` para reduzir pico de RAM.

**Se o build morrer com `signal: killed` (OOM)** — comum em instâncias ~1 GB:

1. Adicione swap no host (exemplo Linux):

```bash
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

2. Ou faça o build noutra máquina com mais RAM e envie a imagem para o servidor (`docker save` / registry).

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

| Secret | Usado por | Conteúdo |
|--------|-----------|----------|
| `postgres_password` | Postgres | Senha do utilizador `campus_connect` |
| `database_url` | API | URL completa (`postgres://...@db:5432/...`) |
| `api_secret` | API | Segredo HMAC dos tokens JWT (string longa aleatória) |

```bash
printf "SENHA_FORTE_AQUI" | docker secret create postgres_password -
printf "postgres://campus_connect:Veiga.2004x@db:5432/campus_connect?sslmode=disable" | docker secret create database_url -
openssl rand -base64 48 | docker secret create api_secret -
```

**Não** coloque senhas nem `API_SECRET` em `.env.db` / `.env.api` — esses ficheiros podem ir para o repositório como exemplos de tuning.

## 4) Preparar variaveis

```bash
cp .env.db.example .env.db
cp .env.api.example .env.api
```

## 5) Deploy das stacks

O `docker stack deploy` **nao suporta** `--env-file` em versoes antigas do Docker CLI.
Use o script `deploy.sh`, que carrega o `.env` no shell e substitui `${VAR}` nos YAML:

```bash
chmod +x deploy.sh

# Banco primeiro (cria a rede overlay campus_overlay_shared)
./deploy.sh campus-db db-stack.yml .env.db

# API depois
./deploy.sh campus-api api-stack.yml .env.api
```

Alternativa manual:

```bash
set -a && . ./.env.db && set +a
docker stack deploy -c db-stack.yml campus-db

set -a && . ./.env.api && set +a
docker stack deploy -c api-stack.yml campus-api
```

## 6) Verificacao

```bash
docker stack services campus-db
docker stack services campus-api
docker service logs -f campus-api_api
```

## Observacoes

- O host do banco dentro da rede overlay e `db`.
- A API lê `DATABASE_URL_FILE` e `API_SECRET_FILE` em `/run/secrets/`.
- `CORS_ORIGIN` e `PUBLIC_BASE_URL` podem ficar em `environment` no stack (não são secret).
- Secret no Swarm nao e editavel; remova e recrie quando precisar trocar valor.
