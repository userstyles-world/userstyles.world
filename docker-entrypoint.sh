#!/bin/sh
set -Eeo pipefail
set -x
[ -f "${DATA_DIR}/${DB}" ] || DB_DROP=1 DB_MIGRATE=1 userstyles
exec "$@"
