# Project 2 implementation notes

## Data model and recovery

`serviceContext` stores JSON provider lists under `nsf/service/<service>` keys.
Each provider contains its NF name, SBI address and SBI port. The database is
the source of truth: registration is saved before success is returned, discovery
reads persisted data, and shutdown never deletes registrations. This also allows
recovery after an abrupt process exit without a separate cache rebuild.

A mutex serializes read-modify-write mutations within one NSF process. Discovery
uses a read lock. Multiple providers are preserved; discovery chooses the first
registered provider, and deregistration removes only the named NF from the named
services. This selection policy is an implementation choice, not an added course
requirement. Multiple independent NSF writers sharing one DB are out of scope.

The existing SBI and processor layers and `dbContext.go` are unchanged.

## MongoDB backend

`internal/db/customDb.go` implements every `DbIf` method with the official MongoDB
Go driver. Document `_id` is the DB key; `value` contains the serialized provider
list. Operations use five-second contexts. `SaveAll` uses ordered bulk upserts;
MongoDB bulk writes are not a multi-document transaction.

For a local demonstration:

```bash
docker run -d --name nsf-pj2-mongo -p 127.0.0.1:27018:27017 mongo:8.0
export NSF_MONGO_URI=mongodb://127.0.0.1:27018
export NSF_MONGO_DATABASE=nsf_project2
```

Set `nsf.db_type` to `mongodb` in a separate demonstration configuration and run
`./build/nsf -c <demo-config.yaml>`. The committed `config.yaml` stays `bbolt`.
Use `docker exec nsf-pj2-mongo mongosh --quiet nsf_project2 --eval
'db.services.find().toArray()'` to inspect persisted registrations.

## Verification

```bash
make lint
make
make test
go test -race -count=1 ./internal/nsfContext
NSF_TEST_MONGO=1 NSF_MONGO_DATABASE=nsf_contract_test go test -count=1 -v ./internal/db
docker run --rm -v "$(pwd):/work" alonza0314/nsf-tester:latest
```

MongoDB tests are opt-in so the default build/test does not require a running
MongoDB server. They run the same existing database contract checks against the
second implementation, without removing or weakening bbolt tests.

Additional real HTTP checks, with an isolated bbolt data path:

```bash
go build -ldflags '-X nsf/internal/db.NSF_DB_PATH=/tmp/nsf-pj2-api.db' -o build/nsf-integration main.go
python3 scripts/test_api.py --db bbolt
NSF_MONGO_DATABASE=nsf_api_test python3 scripts/test_api.py --db mongodb
```

The additional checks cover group discovery, partial deregistration, shared
services, unknown services, graceful restart and abrupt process termination.
They are supplementary checks, not the TA's hidden tests.
