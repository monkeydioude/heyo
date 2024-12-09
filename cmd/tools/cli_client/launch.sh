#!/bin/bash

# Function to generate a random substring of a specified length
generate_random_substring() {
  local length=$1
  # Use base64 to ensure valid UTF-8 characters
  head -c "$((length * 3 / 4 + 1))" < /dev/urandom | base64 | tr -dc 'a-zA-Z0-9' | head -c "$length"
}

# Function to generate the final random string
generate_string() {
  local num_substrings=$1
  local final_string=""
  
  for ((i = 0; i < num_substrings; i++)); do
    # Generate a random length for each substring (e.g., 4 to 12 characters)
    local substring_length=$((RANDOM % 9 + 4)) # Random length between 4 and 12
    local random_substring
    random_substring=$(generate_random_substring "$substring_length")
    final_string+=" $random_substring"
  done

  echo "$final_string"
}

# Function to send a message
sendMessage() {
  # Generate a random number of substrings (e.g., 3 to 10 substrings)
#   local num_substrings=$((RANDOM % 8 + 3)) # Random number between 3 and 10
  local num_substrings=100 # Random number between 3 and 10
  
  # Generate the random string
  local final_string
  final_string=$(generate_string "$num_substrings")

  # Execute the go command with a valid UTF-8 string
  echo "go run main.go -event test1 $final_string" 
  go run main.go -event test1 $final_string &
}

# Main loop
for k in {1..10}; do
  sendMessage
  sendMessage
  sendMessage
  sendMessage
  sendMessage
  echo "Iteration: $k"
done

# Wait for all background processes to finish
wait
