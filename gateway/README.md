### Authorization Key
```
Authorization 19766804-2d26-4c02-ba66-751cada5cbbc
```

### Port
PORT: 11000

### Swagger
https://0.0.0.0:11000/

### Create a payment
```sh
curl --location --request POST 'https://localhost:11000/api/v1/payment' \
--header 'Authorization: 19766804-2d26-4c02-ba66-751cada5cbbc' \
--header 'X-Idempotency-Key: b79fd97c-9aa0-4fa9-84c5-f9667b02b5a0' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 388,
    "currency": "USD",
    "reference": "37447a8a-5ade-441f-a9d1-d49741bcd0d1",
    "credit_card": {
        "first_name": "Maiya",
        "last_name": "Becker",
        "number": "4485040371536584",
        "expiry_month": "12",
        "expiry_year": "2022"
    }
}' --insecure 
```

### Get a payment by id
```sh
curl --location --request GET 'https://localhost:11000/api/v1/payment/11a349b5-8ec3-4619-b4be-8246d824b7ce' \
--header 'Authorization: 19766804-2d26-4c02-ba66-751cada5cbbc' --insecure
```