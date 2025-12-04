.PHONY: env deps whisper
SHELL := /bin/bash
.ONESHELL:

whisper:
	source ./scripts/setup-whisper.sh

deps:
	./scripts/install.sh

env:
	awk -F'=' 'BEGIN {OFS="="} \
    	/^[[:space:]]*#/ {print; next} \
    	/^[[:space:]]*$$/ {print ""; next} \
    	NF>=1 {gsub(/^[[:space:]]+|[[:space:]]+$$/, "", $$1); print $$1"="}' .env > .env.example
	echo ".env.example generated successfully."