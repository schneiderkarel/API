#!/usr/bin/env bash

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

((REPEAT = $2))
while [[ ${REPEAT} -ne 0 ]] ; do
	CODE=$(docker exec $1 sh -c 'pg_isready -q -t 0' &> /dev/null) || CODE=1

	if [[ ${CODE} -eq 0 ]] ; then
		echo -e ${GREEN}'OK: PostgreSQL running in container "'$1'" ping successful'${NC}
		exit 0
	fi

	((REPEAT = REPEAT - 1))

	sleep 1
done

echo -e ${RED}'ERROR: PostgreSQL running in container "'$1'" ping failed'
exit 1
