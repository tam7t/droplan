#!/bin/sh

set -ex

if [ -z "$DO_INTERVAL" ]; then
	# if the interval is not set, only execute once
	./droplan "$@"
else
	while true; do
		# since we use 'set -e', this while loop will exit if droplan exits with a return value other than 0
		# (which in turn tells docker to restart the container (assuming the --restart option was used)
		# while delaying retries exponentially)
		./droplan "$@"
		sleep "$DO_INTERVAL"
	done
fi

