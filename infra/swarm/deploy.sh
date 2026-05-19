#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")"

stack="${1:?uso: ./deploy.sh <stack> <compose.yml> [env-file]}"
compose="${2:?uso: ./deploy.sh <stack> <compose.yml> [env-file]}"
env_file="${3:-}"

if [ -n "$env_file" ]; then
	if [ ! -f "$env_file" ]; then
		echo "arquivo de ambiente nao encontrado: $env_file" >&2
		exit 1
	fi
	set -a
	# shellcheck disable=SC1090
	. "$env_file"
	set +a
fi

docker stack deploy -c "$compose" "$stack"
echo "deploy: $stack ($compose)"
