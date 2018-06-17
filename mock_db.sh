# This script handles local database creation to make testing easier.
# Run with `make db`

postgres -V
if [ "$?" -gt "0" ]; then
  echo "MAKE: Postgres not installed!"
  exit 1
fi

echo "MAKE: Killing existing postgres processes..."
pg_ctl -D /usr/local/var/postgres stop -s -m fast
pg_ctl -D /usr/local/var/postgres start
while true; do
  pg_ctl -D /usr/local/var/postgres status
  if [ "$?" -gt "0" ]; then
    echo "MAKE: Waiting for postgres to start..."
    sleep 3
  else
    break
  fi
done

createuser -s rocket_test
createdb rocket_test_db
psql -d rocket_test_db -a -f ./schema/tables.sql
echo "MAKE: Local database ready to go!"
