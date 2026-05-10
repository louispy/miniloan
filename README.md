# MiniLoan
Implementation of the billing engine problem — provides loan schedule, outstanding amount, delinquency check, and payment processing.

## Setup
This project is built with `Go` version `25.5+`.
Minimum compability version is `22` (set in go.mod).

Copy `.env.example` file into a new `.env` file in project root directory.
Change database url with your own.

This project uses PostgreSQL as database dependency.
Recommended PostgreSQL version: `14+`

Table schemas can be found in `sql` folder.
Alternatively, for convenience, a setup/ migrate main program is provided.

To run database tables setup with Makefile:
```
make run-migration
```

To run tests:
```
go test -v ./...
```

To run http server with Makefile:
```
make run-http
```

## Project structure overview
```
cmd/
  httpserver/   - HTTP server entry point
  migrate/      - one-shot SQL bootstrap
internal/
  services/     - business logic (the spec methods)
  domain/       - models + repositories (SQL)
  database/     - connection + tx manager
  api/          - HTTP handlers
sql/            - table schemas
```

## Design Decisions
- App is separated into layers: business logics (services), database (repositories), and http (out of task scope).
- On loan creation, pre-create all installment rows so reads are easier and no on-the-fly counts/sums are needed.
- `MakePayment` additionally requires `installment_id` which is not in task scope, however may be beneficial to idempotency.
- In `MakePayment`, `installment_id` can be used as a unique identifier for row locking to prevent duplication.
- Installment N due date = loan creation + 7N days, evaluated start-of-day in the business timezone.
- No foreign-key constraints for schema simplicity.
- Unit tests use mocked function.
- Only services are tested (with unit test) since it's the minimal task scope.
- HTTP API layered on top of the service interface so flows can be reviewed with curl

## Sample Curls

### Create Loan

```
curl --request POST \
  --url http://localhost:8989/loans
```
Sample Response
`installments` array is redacted for docs simplicity

```
{
  "data": {
    "id": "564d5e8e-fd06-45b4-b3c9-69fbcf23093e",
    "installments": [
      {
        "id": "e8e70c6c-893a-45e2-a6f4-8d1ec50539cf",
        "week": 1,
        "amount": 110000,
        "status": 1,
        "due_date": "2026-05-17 00:00:00 +0700"
      }
    ]
  },
  "message": "Successfully created a new loan request"
}
```

### Get Installments

```
curl --request GET \
  --url 'http://localhost:8989/loans/:id/installments?status=unpaid' \
  --header 'content-type: application/json'
```

Sample Response
`installments` array is redacted for docs simplicity

```
{
  "data": {
    "installments": [
      {
        "id": "c1bac486-b0a9-4abb-8eea-994267d56d0f",
        "week": 2,
        "amount": 110000,
        "status": 1,
        "due_date": "2026-05-24 00:00:00 +0000"
      },
    ]
  },
  "message": "Successfully get installments"
}
```

###  GetOutstanding

```
curl --request GET \
  --url http://localhost:8989/loans/:id/outstanding

```

Sample Response

```
{
  "data": {
    "amount": 5500000
  },
  "message": "Successfully get outstanding amount"
}
```

### IsDeliquent

```
curl --request GET \
  --url http://localhost:8989/loans/:id/deliquent
```

Sample Response

```
{
  "data": {
    "is_deliquent": false
  },
  "message": "Successfully get IsDeliquent value"
}
```

### MakePayment

```
curl --request POST \
  --url http://localhost:8989/loans/:id/payment \
  --header 'content-type: application/json' \
  --data '{
  "amount": 110000,
  "installment_id": "e8e70c6c-893a-45e2-a6f4-8d1ec50539cf"
}'
```

Sample Response
```
{
  "data": {
    "installment_id": "c1bac486-b0a9-4abb-8eea-994267d56d0f"
  },
  "message": "Successfully MakePayment"
}
```
