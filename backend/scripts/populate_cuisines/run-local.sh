#!/bin/bash

if [ -f .venv ]; then
	source .venv/bin/activate
fi

python scripts/populate_cuisines.py
