#!/bin/bash

# Test script for authentication API
BASE_URL="http://localhost:8081"

echo "==================================="
echo "Testing Authentication System"
echo "==================================="

echo -e "\n1. Register with Email + Password"
curl -X POST "$BASE_URL/auth/register-password" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "password123",
    "role": "student"
  }' | jq '.'

echo -e "\n\n2. Register with Phone + Password"
curl -X POST "$BASE_URL/auth/register-password" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+77001234567",
    "password": "password123",
    "role": "employer"
  }' | jq '.'

echo -e "\n\n3. Login with Email + Password"
curl -X POST "$BASE_URL/auth/login-password" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "student@example.com",
    "password": "password123"
  }' | jq '.'

echo -e "\n\n4. Login with Phone + Password"
curl -X POST "$BASE_URL/auth/login-password" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "+77001234567",
    "password": "password123"
  }' | jq '.'

echo -e "\n\n5. Request Email OTP Code"
curl -X POST "$BASE_URL/auth/request-email-code" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com"
  }' | jq '.'

echo -e "\n(Check server logs for OTP code)"
echo -e "Enter the OTP code from server logs:"
read OTP_CODE

echo -e "\n\n6. Verify Email OTP Code"
curl -X POST "$BASE_URL/auth/verify-email-code" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"newuser@example.com\",
    \"code\": \"$OTP_CODE\"
  }" | jq '.'

echo -e "\n\n7. Request Phone OTP Code"
curl -X POST "$BASE_URL/auth/request-phone-code" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+77009876543",
    "role": "student"
  }' | jq '.'

echo -e "\n(Check server logs for OTP code)"
echo -e "Enter the OTP code from server logs:"
read PHONE_OTP_CODE

echo -e "\n\n8. Verify Phone OTP Code"
curl -X POST "$BASE_URL/auth/verify-phone-code" \
  -H "Content-Type: application/json" \
  -d "{
    \"phone\": \"+77009876543\",
    \"code\": \"$PHONE_OTP_CODE\"
  }" | jq '.'

echo -e "\n\n==================================="
echo "All tests completed!"
echo "==================================="
