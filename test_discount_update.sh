#!/bin/bash

# Configuration
BASE_URL="https://bishalchhetri.com.np/api/v1"
SALE_ID="81"
CHARGE_ID="259"
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJjb21wYW55X2lkIjoxLCJ1c2VybmFtZSI6InNha3NoeWFtIiwiZXhwIjoxNzkyMzk1MTI4fQ.Lju1xLmmxcgbKuoS-Y83tD2g4vD6cWyuS-OZScPOIEk"

echo "Testing Discount Update API..."
echo "URL: ${BASE_URL}/sales/${SALE_ID}/charges/${CHARGE_ID}"
echo "=========================================="

# Test the discount update API
curl -X PUT "${BASE_URL}/sales/${SALE_ID}/charges/${CHARGE_ID}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${JWT_TOKEN}" \
  -d '{
    "charge_type": "discount",
    "amount": -2000.0,
    "description": "Updated discount amount via API test"
  }' \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -v

echo "=========================================="
echo "Test completed!"
