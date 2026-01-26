#!/bin/sh
set -eu

# Wipes ALL records from the MySQL database (keeps tables).
# Intended to be run inside the `mysql` container:
#   docker compose exec -e CONFIRM=YES mysql sh /scripts/wipe-db.sh
#
# Safety: requires CONFIRM=YES

if [ "${CONFIRM:-}" != "YES" ]; then
  echo "Refusing to run: this will DELETE ALL DATA in the database."
  echo "Re-run with: CONFIRM=YES"
  exit 1
fi

# Prefer MySQL container env vars, fallback to app env vars.
DB_NAME="${MYSQL_DATABASE:-${DATABASE_NAME:-}}"
if [ -z "$DB_NAME" ]; then
  echo "Database name not set. Provide MYSQL_DATABASE or DATABASE_NAME."
  exit 1
fi

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"

# Use root if available (most reliable permissions), otherwise use regular user.
if [ -n "${MYSQL_ROOT_PASSWORD:-}" ]; then
  MYSQL_USER="root"
  MYSQL_PWD="${MYSQL_ROOT_PASSWORD}"
else
  MYSQL_USER="${MYSQL_USER:-${DATABASE_USER:-}}"
  MYSQL_PWD="${MYSQL_PASSWORD:-${DATABASE_PASSWORD:-}}"
fi

if [ -z "${MYSQL_USER:-}" ] || [ -z "${MYSQL_PWD:-}" ]; then
  echo "MySQL credentials not set. Provide MYSQL_ROOT_PASSWORD or MYSQL_USER + MYSQL_PASSWORD."
  exit 1
fi

echo "Wiping all data from MySQL database: ${DB_NAME}"

BASE_ARGS="-h${MYSQL_HOST} -P${MYSQL_PORT} -u${MYSQL_USER}"

TABLES_FILE="$(mktemp)"
AUTOINC_FILE="$(mktemp)"
cleanup() {
  rm -f "$TABLES_FILE" "$AUTOINC_FILE"
}
trap cleanup EXIT

# List base tables in the schema
MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS -N -B -e \
  "SELECT table_name
   FROM information_schema.tables
   WHERE table_schema='${DB_NAME}'
     AND table_type='BASE TABLE';" > "$TABLES_FILE"

if [ ! -s "$TABLES_FILE" ]; then
  echo "No tables found in schema '${DB_NAME}'. Nothing to wipe."
  exit 0
fi

# List tables that have AUTO_INCREMENT (so we can safely reset counters)
MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS -N -B -e \
  "SELECT table_name
   FROM information_schema.tables
   WHERE table_schema='${DB_NAME}'
     AND AUTO_INCREMENT IS NOT NULL;" > "$AUTOINC_FILE"

{
  echo "SET FOREIGN_KEY_CHECKS=0;"

  while IFS= read -r t; do
    [ -n "$t" ] && echo "DELETE FROM \`$t\`;"
  done < "$TABLES_FILE"

  while IFS= read -r t; do
    [ -n "$t" ] && echo "ALTER TABLE \`$t\` AUTO_INCREMENT = 1;"
  done < "$AUTOINC_FILE"

  echo "SET FOREIGN_KEY_CHECKS=1;"
} | MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS "$DB_NAME"

echo "Done."
#!/bin/sh
set -eu

# Wipes ALL records from the MySQL database (keeps tables).
# Intended to be run inside the `mysql` container:
#   docker compose exec mysql sh /scripts/wipe-db.sh
#
# Safety: requires CONFIRM=YES
#   docker compose exec -e CONFIRM=YES mysql sh /scripts/wipe-db.sh

if [ "${CONFIRM:-}" != "YES" ]; then
  echo "Refusing to run: this will DELETE ALL DATA in the database."
  echo "Re-run with: CONFIRM=YES"
  exit 1
fi

# Prefer MySQL container env vars, fallback to app env vars.
DB_NAME="${MYSQL_DATABASE:-${DATABASE_NAME:-}}"
if [ -z "$DB_NAME" ]; then
  echo "Database name not set. Provide MYSQL_DATABASE or DATABASE_NAME."
  exit 1
fi

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"

# Use root if available (most reliable permissions), otherwise use regular user.
if [ -n "${MYSQL_ROOT_PASSWORD:-}" ]; then
  MYSQL_USER="root"
  MYSQL_PWD="${MYSQL_ROOT_PASSWORD}"
else
  MYSQL_USER="${MYSQL_USER:-${DATABASE_USER:-}}"
  MYSQL_PWD="${MYSQL_PASSWORD:-${DATABASE_PASSWORD:-}}"
fi

if [ -z "${MYSQL_USER:-}" ] || [ -z "${MYSQL_PWD:-}" ]; then
  echo "MySQL credentials not set. Provide MYSQL_ROOT_PASSWORD or MYSQL_USER + MYSQL_PASSWORD."
  exit 1
fi

echo "Wiping all data from MySQL database: ${DB_NAME}"

BASE_ARGS="-h${MYSQL_HOST} -P${MYSQL_PORT} -u${MYSQL_USER}"

TABLES_FILE="$(mktemp)"
AUTOINC_FILE="$(mktemp)"
cleanup() {
  rm -f "$TABLES_FILE" "$AUTOINC_FILE"
}
trap cleanup EXIT

# List base tables in the schema
MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS -N -B -e \
  "SELECT table_name
   FROM information_schema.tables
   WHERE table_schema='${DB_NAME}'
     AND table_type='BASE TABLE';" > "$TABLES_FILE"

if [ ! -s "$TABLES_FILE" ]; then
  echo "No tables found in schema '${DB_NAME}'. Nothing to wipe."
  exit 0
fi

# List tables that have AUTO_INCREMENT (so we can safely reset counters)
MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS -N -B -e \
  "SELECT table_name
   FROM information_schema.tables
   WHERE table_schema='${DB_NAME}'
     AND AUTO_INCREMENT IS NOT NULL;" > "$AUTOINC_FILE"

{
  echo "SET FOREIGN_KEY_CHECKS=0;"

  while IFS= read -r t; do
    [ -n "$t" ] && echo "DELETE FROM \`$t\`;"
  done < "$TABLES_FILE"

  while IFS= read -r t; do
    [ -n "$t" ] && echo "ALTER TABLE \`$t\` AUTO_INCREMENT = 1;"
  done < "$AUTOINC_FILE"

  echo "SET FOREIGN_KEY_CHECKS=1;"
} | MYSQL_PWD="$MYSQL_PWD" mysql $BASE_ARGS "$DB_NAME"

echo "Done."
