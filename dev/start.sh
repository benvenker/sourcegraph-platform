#!/bin/bash

if [ -n "$DELVE_FRONTEND" ]; then
	export DELVE=1
	echo 'Launching frontend with delve'
	export EXEC_FRONTEND='dlv exec --headless --listen=:2345 --log'
fi

if [ -n "$DELVE_SEARCHER" ]; then
	export DELVE=1
	echo 'Launching searcher with delve'
	export EXEC_SEARCHER='dlv exec --headless --listen=:2346 --log'
fi

if [ -n "$DELVE" ]; then
	echo 'Due to a limitation in delve, bebug binaries will not start until you attach a debugger.'
	echo 'See https://github.com/derekparker/delve/issues/952'
fi

set -euf -o pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.." # cd to repo root dir

export AUTH0_CLIENT_ID=onW9hT0c7biVUqqNNuggQtMLvxUWHWRC
export AUTH0_CLIENT_SECRET=cpse5jYzcduFkQY79eDYXSwI6xVUO0bIvc4BP6WpojdSiEEG6MwGrt8hj_uX3p5a
export AUTH0_DOMAIN=sourcegraph-dev.auth0.com
export AUTH0_MANAGEMENT_API_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJSYW1KekRwRmN6SFZZNTBpcmFSb0JMdTNRVmFHTE1VRiIsInNjb3BlcyI6eyJ1c2VycyI6eyJhY3Rpb25zIjpbInJlYWQiLCJ1cGRhdGUiXX0sInVzZXJfaWRwX3Rva2VucyI6eyJhY3Rpb25zIjpbInJlYWQiXX0sInVzZXJzX2FwcF9tZXRhZGF0YSI6eyJhY3Rpb25zIjpbInVwZGF0ZSJdfX0sImlhdCI6MTQ3NzA5NDQxOSwianRpIjoiMTA3YzYyMTZjNWZjYzVjNGNkYjYzZTgxNjRjYjg3ODgifQ.ANOcIGeFPH7X_ppl-AXcv2m0zI7hWwqDlRwJ6h_rMdI
export GITHUB_CLIENT_ID=6f2a43bd8877ff5fd1d5
export GITHUB_CLIENT_SECRET=c5ff37d80e3736924cbbdf2922a50cac31963e43
export LIGHTSTEP_PROJECT=sourcegraph-dev
export LIGHTSTEP_ACCESS_TOKEN=d60b0b2477a7ccb05d7783917f648816
export LIGHTSTEP_INCLUDE_SENSITIVE=true
export PGSSLMODE=disable
export PUBLIC_REPO_REDIRECTS=false
export AUTO_REPO_ADD=true

export SRC_APP_SECRET_KEY=OVSHB1Yru3rlsQ0eKNi2GXCZ47zU7DCK
export GITHUB_BASE_URL=http://127.0.0.1:3180
export SRC_REPOS_DIR=$HOME/.sourcegraph/repos
export DEBUG=true
export SRC_APP_DISABLE_SUPPORT_SERVICES=true
export SRC_GIT_SERVERS=127.0.0.1:3178
export SEARCHER_URL=http://127.0.0.1:3181
export LSP_PROXY=127.0.0.1:4388
export REDIS_MASTER_ENDPOINT=127.0.0.1:6379
export SRC_SESSION_STORE_REDIS=127.0.0.1:6379
export SRC_INDEXER=127.0.0.1:3179
export SRC_SYNTECT_SERVER=http://localhost:3700
export SRC_FRONTEND_INTERNAL=localhost:3090

export REPO_LIST_UPDATE_INTERVAL=1
export PHABRICATOR_URL="http://phabricator.sgdev.org"
export CORS_ORIGIN="https://github.com http://phabricator.sgdev.org"

# Uncomment to sync gitolite repositories.
# export GITOLITE_HOSTS="gitolite.sgdev.org/!git@gitolite.sgdev.org"

# Uncomment one or the other of these configs depending on your testing needs.
# Both "token" and "repos" are optional but you need to have at least one set.
#
# export PHABRICATOR_CONFIG='[{"url":"http://phabricator.sgdev.org","token":"api-agswx2nwodkweitoo3t5l4dcc5xu"}]'
export PHABRICATOR_CONFIG='[{"url":"http://phabricator.sgdev.org","repos":[{"path":"gitolite.sgdev.org/testing","callsign":"TEST"}]}]'

export GITHUB_CONFIG='[{"url": "https://ghe.sgdev.org", "token":"23993bbf8e0fee068b8f70db05fc445d5a7a83da"}]'

export LANGSERVER_GO=${LANGSERVER_GO-"tcp://localhost:4389"}
export LANGSERVER_GO_BG=${LANGSERVER_GO_BG-"tcp://localhost:4389"}

export LICENSE_KEY=${LICENSE_KEY:-24348deeb9916a070914b5617a9a4e2c7bec0d313ca6ae11545ef034c7138d4d8710cddac80980b00426fb44830263268f028c9735}

if ! [ -z "${ZOEKT-}" ]; then
	export ZOEKT_HOST=localhost:6070
fi

# WebApp
export NODE_ENV=development

# Make sure chokidar-cli is installed
npm install

./dev/go-install.sh

# Increase ulimit (not needed on Windows/WSL)
type ulimit > /dev/null && ulimit -n 10000 || true

export GOREMAN=".bin/goreman -f dev/Procfile"
exec $GOREMAN start
