#!/bin/bash

if [ -f .env ]; then
	source .env
fi

cd sql/schema || exit
goose postgres "$DATABASE_URL" down
