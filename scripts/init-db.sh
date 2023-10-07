#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE proddatabase;

	CREATE TABLE IF NOT EXISTS avatars (
	id SERIAL PRIMARY KEY,
	path VARCHAR
	);
EOSQL