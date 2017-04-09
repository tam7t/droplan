#!/bin/sh

set -ex

if [ -z "$DO_INTERVAL" ]; then
	# specifies the check interval
	DO_INTERVAL=300
fi

while true; do
	# since we use 'set -e', this while loop will exit if droplan exits with a return value other than 0
	# (which in turn tells docker to restart the container (assuming the --restart option was used)
	# while delaying retries exponentially)
	./droplan "$@"
	sleep "$DO_INTERVAL"
done
