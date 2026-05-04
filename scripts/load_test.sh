#!/bin/bash

# Load Simulation Script for Assignment 6
# Simulates concurrent requests to backend services to test capacity

URL=${1:-"http://localhost/api/orders"}
CONCURRENT_REQUESTS=${2:-10}
TOTAL_REQUESTS=${3:-100}

echo "Starting load test on $URL"
echo "Concurrency: $CONCURRENT_REQUESTS, Total requests: $TOTAL_REQUESTS"

# Check if 'ab' (Apache Benchmark) is installed, otherwise use a simple loop
if command -v ab &> /dev/null
then
    ab -n $TOTAL_REQUESTS -c $CONCURRENT_REQUESTS $URL/health
else
    echo "Apache Benchmark (ab) not found, using curl in a loop..."
    for i in $(seq 1 $TOTAL_REQUESTS); do
        curl -s -o /dev/null -w "%{http_code}\n" $URL/health &
        if (( $i % $CONCURRENT_REQUESTS == 0 )); then
            wait
        fi
    done
    wait
fi

echo "Load test completed."
