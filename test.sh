#!/bin/bash

BASE_URL="http://13.127.151.58:2000"  # Replace with your actual base URL
NUM_REQUESTS=20

# Function to make a single request to /api/v1/hello
make_hello_request() {
    start_time=$(date +%s.%N)
    # Capture the response body in one file and the status code separately
    response=$(curl -s -w "%{http_code}" -o temp_response.txt $BASE_URL/api/v1/hello)
    end_time=$(date +%s.%N)
    
    # Calculate duration in milliseconds
    duration=$(echo "($end_time - $start_time) * 1000" | bc)
    
    # Check if status code is an integer and handle errors
    if ! [[ "$response" =~ ^[0-9]+$ ]]; then
        response="000"  # Use a default status code if response isn't numeric (like "Internal Error")
    fi
    echo "$response $duration"
}

# Make requests to /api/v1/hello
echo "Making $NUM_REQUESTS requests to /api/v1/hello..."

successful_requests=0
failed_requests=0
total_time=0

for i in $(seq 1 $NUM_REQUESTS); do
    result=$(make_hello_request)
    status_code=$(echo $result | cut -d' ' -f1)
    response_time=$(echo $result | cut -d' ' -f2)

    echo "Request $i:"
    if [ "$status_code" -eq 200 ]; then
        echo "Response body: $(cat temp_response.txt)"
        successful_requests=$((successful_requests + 1))
    else
        # If the status code isn't 200, display the error message
        echo "Failed with status code: $status_code"
        error_response=$(cat temp_response.txt)

        if [ -n "$error_response" ]; then
            echo "Error response: $error_response"
        else
            echo "No response body (possible network error)"
        fi
        failed_requests=$((failed_requests + 1))
    fi

    # Accumulate total time for calculating the average
    total_time=$(echo "$total_time + $response_time" | bc)
    echo
done

total_requests=$NUM_REQUESTS

# Print hello endpoint stats
echo
echo "Hello Endpoint Stats:"
echo "{
  \"successful_requests\": $successful_requests,
  \"failed_requests\": $failed_requests,
  \"total_requests\": $total_requests,
}"

# Get worker stats
echo
echo "Worker Stats:"
curl -s $BASE_URL/worker/stats | jq '.'

# Cleanup temp file
rm temp_response.txt

