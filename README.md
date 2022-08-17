### To check service work:

1. **To check service works locally please add** _Undefined_ to the allowedCountiesCodes in config.yaml

2. **Start all services (database and queue)**

```
docker-compose up --build -d app pgweb
```

3. **Create new company. Please create several companies**

```
curl --location --request POST 'http://127.0.0.1:8080/api/v1/company' \
    --header "Content-Type: application/json" \
    --data-raw '{
          "name": "company_name",
          "code":"company_code",
          "country": "United Kindom",
          "website": "https://example.com",
          "phone": "+1401235566"
      }'
```

```
curl --location --request POST 'http://127.0.0.1:8080/api/v1/company' \
    --header "Content-Type: application/json" \
    --data-raw '{
          "name": "11",
          "code":"22",
          "country": "United Kindom",
          "website": "https://example.com",
          "phone": "+1401235566"
      }'
```

```
curl --location --request POST 'http://127.0.0.1:8080/api/v1/company' \
    --header "Content-Type: application/json" \
    --data-raw '{
          "name": "111",
          "code":"222",
          "country": "France",
          "website": "https://example.com",
          "phone": "+1401235566"
      }'
```

```
curl --location --request POST 'http://127.0.0.1:8080/api/v1/company' \
    --header "Content-Type: application/json" \
    --data-raw '{
          "name": "11",
          "code":"2223",
          "country": "Ukraine",
          "website": "https://example.com",
          "phone": "+1401235566"
      }'
```

4. **We can get some company by his name and code**

```
curl --location --request GET 'http://127.0.0.1:8080/api/v1/company/company_name/company_code'
```

```
curl --location --request GET 'http://127.0.0.1:8080/api/v1/company/11/22'
```

5. **We can filter some companies**

```
curl --location --request GET 'http://127.0.0.1:8080/api/v1/companies?limit=3&name=1&code=3'
```

```
curl --location --request GET 'http://127.0.0.1:8080/api/v1/companies?limit=3&name=1&code=2'
```

6. **Update some company**

```
curl --location --request PUT 'http://127.0.0.1:8080/api/v1/company/11/22' \
    --header "Content-Type: application/json" \
    --data-raw '{
          "name": "112",
          "code":"2223",
          "country": "Germany",
          "website": "https://example.com",
          "phone": "+1401235566"
      }'
```

### To create first migration schema please use this command:

```
 docker run -v $(pwd)/migration:/migrations migrate/migrate create -ext sql -dir=/migrations -seq create_companies_tbl
```
