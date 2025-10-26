#!/bin/bash

# String Analyzer API Test Script
# Usage: ./test_api.sh [base_url]
# Example: ./test_api.sh http://localhost:8080

BASE_URL="${1:-http://localhost:8080}"

echo "========================================="
echo "String Analyzer API Test Suite"
echo "========================================="
echo "Testing API at: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

test_count=0
pass_count=0
fail_count=0

# Function to test endpoint
test_endpoint() {
    test_count=$((test_count + 1))
    local description=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    
    echo -e "${BLUE}Test $test_count: $description${NC}"
    echo "  Method: $method"
    echo "  Endpoint: $endpoint"
    
    if [ -n "$data" ]; then
        echo "  Data: $data"
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" == "$expected_status" ]; then
        echo -e "  ${GREEN}✓ PASS${NC} (Status: $http_code)"
        pass_count=$((pass_count + 1))
    else
        echo -e "  ${RED}✗ FAIL${NC} (Expected: $expected_status, Got: $http_code)"
        fail_count=$((fail_count + 1))
    fi
    
    echo "  Response: $body"
    echo ""
}

echo "========================================="
echo "1. HEALTH CHECK"
echo "========================================="
test_endpoint \
    "Health check endpoint" \
    "GET" \
    "/health" \
    "" \
    "200"

echo "========================================="
echo "2. CREATE STRINGS (POST /strings)"
echo "========================================="

test_endpoint \
    "Create palindrome string 'racecar'" \
    "POST" \
    "/strings" \
    '{"value": "racecar"}' \
    "201"

test_endpoint \
    "Create normal string 'hello world'" \
    "POST" \
    "/strings" \
    '{"value": "hello world"}' \
    "201"

test_endpoint \
    "Create single word 'noon'" \
    "POST" \
    "/strings" \
    '{"value": "noon"}' \
    "201"

test_endpoint \
    "Create multi-word palindrome 'A man a plan a canal Panama'" \
    "POST" \
    "/strings" \
    '{"value": "A man a plan a canal Panama"}' \
    "201"

test_endpoint \
    "Create string 'testing'" \
    "POST" \
    "/strings" \
    '{"value": "testing"}' \
    "201"

test_endpoint \
    "Create duplicate string 'racecar' (should fail)" \
    "POST" \
    "/strings" \
    '{"value": "racecar"}' \
    "409"

test_endpoint \
    "Create with missing value field (should fail)" \
    "POST" \
    "/strings" \
    '{}' \
    "400"

test_endpoint \
    "Create with empty value (should fail)" \
    "POST" \
    "/strings" \
    '{"value": ""}' \
    "400"

echo "========================================="
echo "3. GET SPECIFIC STRING"
echo "========================================="

test_endpoint \
    "Get string 'racecar'" \
    "GET" \
    "/strings/racecar" \
    "" \
    "200"

test_endpoint \
    "Get string 'hello world' (URL encoded)" \
    "GET" \
    "/strings/hello%20world" \
    "" \
    "200"

test_endpoint \
    "Get non-existent string (should fail)" \
    "GET" \
    "/strings/nonexistent" \
    "" \
    "404"

echo "========================================="
echo "4. GET ALL STRINGS WITH FILTERS"
echo "========================================="

test_endpoint \
    "Get all strings (no filters)" \
    "GET" \
    "/strings" \
    "" \
    "200"

test_endpoint \
    "Get palindromes only" \
    "GET" \
    "/strings?is_palindrome=true" \
    "" \
    "200"

test_endpoint \
    "Get non-palindromes" \
    "GET" \
    "/strings?is_palindrome=false" \
    "" \
    "200"

test_endpoint \
    "Get single word strings" \
    "GET" \
    "/strings?word_count=1" \
    "" \
    "200"

test_endpoint \
    "Get strings with min_length=5" \
    "GET" \
    "/strings?min_length=5" \
    "" \
    "200"

test_endpoint \
    "Get strings with max_length=10" \
    "GET" \
    "/strings?max_length=10" \
    "" \
    "200"

test_endpoint \
    "Get strings with min and max length" \
    "GET" \
    "/strings?min_length=4&max_length=10" \
    "" \
    "200"

test_endpoint \
    "Get strings containing 'o'" \
    "GET" \
    "/strings?contains_character=o" \
    "" \
    "200"

test_endpoint \
    "Combined filters: palindrome + single word" \
    "GET" \
    "/strings?is_palindrome=true&word_count=1" \
    "" \
    "200"

echo "========================================="
echo "5. NATURAL LANGUAGE FILTERING"
echo "========================================="

test_endpoint \
    "NL Query: single word palindromes" \
    "GET" \
    "/strings/filter-by-natural-language?query=single%20word%20palindromes" \
    "" \
    "200"

test_endpoint \
    "NL Query: strings longer than 10 characters" \
    "GET" \
    "/strings/filter-by-natural-language?query=strings%20longer%20than%2010%20characters" \
    "" \
    "200"

test_endpoint \
    "NL Query: containing letter o" \
    "GET" \
    "/strings/filter-by-natural-language?query=containing%20letter%20o" \
    "" \
    "200"

test_endpoint \
    "NL Query: palindromic strings with first vowel" \
    "GET" \
    "/strings/filter-by-natural-language?query=palindromic%20strings%20with%20first%20vowel" \
    "" \
    "200"

test_endpoint \
    "NL Query: missing query parameter (should fail)" \
    "GET" \
    "/strings/filter-by-natural-language" \
    "" \
    "400"

echo "========================================="
echo "6. DELETE STRINGS"
echo "========================================="

test_endpoint \
    "Delete string 'testing'" \
    "DELETE" \
    "/strings/testing" \
    "" \
    "204"

test_endpoint \
    "Delete same string again (should fail)" \
    "DELETE" \
    "/strings/testing" \
    "" \
    "404"

test_endpoint \
    "Delete non-existent string (should fail)" \
    "DELETE" \
    "/strings/doesnotexist" \
    "" \
    "404"

echo "========================================="
echo "7. VERIFY DELETION"
echo "========================================="

test_endpoint \
    "Get deleted string 'testing' (should fail)" \
    "GET" \
    "/strings/testing" \
    "" \
    "404"

test_endpoint \
    "Get all strings (should not include 'testing')" \
    "GET" \
    "/strings" \
    "" \
    "200"

echo "========================================="
echo "TEST SUMMARY"
echo "========================================="
echo -e "Total Tests: $test_count"
echo -e "${GREEN}Passed: $pass_count${NC}"
echo -e "${RED}Failed: $fail_count${NC}"
echo ""

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi