#!/usr/bin/env bash

set -euo pipefail

exec CompileDaemon -build="go build -o ./cmd/main ./cmd/." -command="./cmd/main"
