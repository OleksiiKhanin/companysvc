### To check service work:

1. **To check service works locally please add** _Undefined_ to the allowedCountiesCodes in config.yaml

2. **Start all services (database and queue)**

```
docker-compose up --build -d app pgweb
```

To check that all services are running you can use this command:
```
docker-compose ps
```

You should see something like this

```
    Name                  Command               State                              Ports
--------------------------------------------------------------------------------------------------------------------
company_svc    /bin/sh -c "companysvc"          Up      0.0.0.0:8080->8080/tcp,:::8080->8080/tcp
pgweb_svc      /usr/bin/pgweb --bind=0.0. ...   Up      0.0.0.0:6080->8081/tcp,:::6080->8081/tcp
postgres_svc   docker-entrypoint.sh postgres    Up      0.0.0.0:5432->5432/tcp,:::5432->5432/tcp
queue_svc      docker-entrypoint.sh nats- ...   Up      0.0.0.0:4222->4222/tcp,:::4222->4222/tcp, 6222/tcp, 8222/tcp

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

### To run tests please execute this command in the project root directory

```
go test -v ./...
```

### TODO
1. Add a JWT authentication feature.
2. Add swagger documentation
