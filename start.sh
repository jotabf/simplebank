#!/bin/sh

# This line sets the shell option "errexit" which causes the shell to 
# exit immediately if a command exits with a non-zero status. This is 
# useful for ensuring that errors are caught and dealt with quickly.
set -e 

echo "Run db migration: $DB_SOURCE"
/app/migrate -path /app/db/migration -database $DB_SOURCE -verbose up
#/app/migrate -path /app/db/migration -database postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable -verbose up

echo "Start the app"

# Execute the command line arguments passed to the script.
# This allows the script to be used as a wrapper for other commands.
exec "$@"