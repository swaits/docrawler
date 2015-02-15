#!/bin/sh

# colorize the output of "go test"

sed ''/ok/s//$(printf "\033[32;1m&\033[0m")/''       | \
	sed ''/PASS/s//$(printf "\033[32;1m&\033[0m")/''   | \
	sed ''/FAIL/s//$(printf "\033[31;1m&\033[0m")/''
