## Bank Simulator
API to simulate an acquirer. It has only one endpoint called transactions of type POST.

- It does not currently validate CVV.
- X-Auth-token is required, it can be any string value.
- Always runs on port 80

### Example
```sh
curl --location --request POST 'http://localhost:80/transactions' \
--header 'X-Auth-token: 4f4baa35-6efe-4734-8220-6fb867ef6959' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 1000,
    "currency": "USD",
    "expiry_month": "10",
    "expiry_year": "2022",
    "first_name": "Amani",
    "last_name": "Lakin",
    "number": "5558468902774508"
}'
```

## Test cards
Authentication successful
---
- 4485040371536584
- 4543474002249996
- 5588686116426417
- 5436031030606378
- 5199992312641465
- 345678901234564
	
Not authenticated
---
- 4539628347117863
- 5309961755464047

Authentication could not be performed
---
- 4024007186645015
- 5234106378657904
	
Attempted authentication
---
- 4556574722325580
- 5558468902774508

Authentication rejected
---
- 4275765574319271
- 5596061690670931
	
Card not enrolled
---
- 4484070000035519
- 5352151570003404
	
Error message during scheme communication
---
- 4452927588210665