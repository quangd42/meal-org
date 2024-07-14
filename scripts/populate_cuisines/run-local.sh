#!/bin/bash

# Run from the root dir
cd scripts/populate_cuisines || exit
go build -o bin/populate_cuisines populate_cuisines.go && ./bin/populate_cuisines
