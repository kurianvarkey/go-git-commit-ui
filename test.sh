#!/bin/sh

echo "Running setup go test in the tests/ folder."

cd ./tests

#Running the go test
#go test -v ./... -count=1 | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
#go test -cover -coverpkg "../src/..." "./..."

go test -v -coverprofile=./coverage/coverage.out ./... -count=1 -cover -coverpkg "../src/..." "./..." | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
