#!/bin/bash

# ANSI color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Find the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/.."

# Variables to track test results
total_tests=0
passed_tests=0
failed_tests=0

# Function to run tests for a specific package
run_test() {
    local package=$1
    echo -e "${YELLOW}Running tests for $package...${NC}"
    
    # Change directory and run tests
    cd "$PROJECT_ROOT"
    go test "./$package" -v
    local test_result=$?
    
    total_tests=$((total_tests + 1))
    if [ $test_result -eq 0 ]; then
        passed_tests=$((passed_tests + 1))
        echo -e "${GREEN}✓ Tests passed for $package${NC}"
    else
        failed_tests=$((failed_tests + 1))
        echo -e "${RED}✗ Tests failed for $package${NC}"
    fi
}

# List of test packages
test_packages=(
    "internal/config"
    "internal/handler"
    "internal/redis"
    "internal/service"
    "pkg/analytics"
    "pkg/errors"
    "pkg/ratelimiter"
    "pkg/shortener"
    "pkg/validator"
)

# Pre-test information
echo -e "${YELLOW}Starting comprehensive test suite...${NC}"
start_time=$(date +%s)

# Run tests for each package
for package in "${test_packages[@]}"; do
    run_test "$package"
done

# Test result summary
end_time=$(date +%s)
duration=$((end_time - start_time))

echo -e "\n${YELLOW}Test Suite Summary:${NC}"
echo -e "Total Tests:  ${total_tests}"
echo -e "${GREEN}Passed Tests: ${passed_tests}${NC}"
echo -e "${RED}Failed Tests: ${failed_tests}${NC}"
echo -e "Total Duration: ${duration} seconds"

# Exit with appropriate status code
if [ $failed_tests -gt 0 ]; then
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
else
    echo -e "${GREEN}All tests passed successfully!${NC}"
    exit 0
fi