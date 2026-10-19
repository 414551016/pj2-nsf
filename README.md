# Network Service Function

This is a network function used for network service register/de-register/discovery.

## SBI Api Endpoint Point

### API Table

| Service | Endpoint | Description | Model |
| - | - | - | - |
| Heartbeat | `/sbi/heartbeat` | health check | - |
| NF Management | `/sbi/nnsf-nfm` | register & de-register network service | [nnsf-nfm](#nnsf-nfm) |
| Service Discovery | `/sbi/nnsf-disc` | discover the service information | [nnsf-disc](#nnsf-disc) |
| Group Service Discovery | `/sbi/nnsf-disc/group` | discover multiple services information | [nnsf-disc](#nnsf-disc) |

### API Structure

```bash
/sbi
    └─/heartbeat(GET)
    └─/nnsf-nfm(POST\DELETE)
    └─/nnsf-disc(Get)
        └─/group(Get)
```

## Make Command

| Command | Description |
| - | - |
| `make` | default set to `make bin` |
| `make bin` | build the executable binary file and output under **build** folder |
| `make clean` | clean up the build folder |
| `make lint` | run golang-lint |
| `make test` | run go test through the whole repo |
| `make tidy` | run module tidy |
| `make run` | run the NSF binary file under **build** folder with [config.yaml](./config.yaml) |
| `make cleandb` | remove the bbolt db file |
| `make cleantestdb` | remove the test bbolt db files |

## What you need TODO?

There are `TODO` keywords at sections that you need to implement. So, you can try to search the `TODO` keywords.

1. `internal/nsfContext/serviceContext.go`: This file is the core context structure for manage the NF services. You have to implement it correctly which can pass the unit test in `serviceContext_test.go` and the integration test with the docker image provided by CNDI.
2. `internal/db/customDb.go`: You can implement your customized `dbIf` at this file and use this DB to test your service is stored correctly.
3. Use the `make` related command to make sure everything is correct.
4. Use the docker image provided by CNDI to validate if your NSF is implement correctly.

## Docker Image Validate

- Pull the image

    ```bash
    docker pull alonza0314/nsf-tester:latest
    ```

- Run the image for testing with local directory

    ```bash
    docker run --rm -v <absolute path to your local directory>:/work alonza0314/nsf-tester:latest
    ```

## SBI Model

### nnsf-nfm

- POST Request

    ```go
    type NnsfNfmPostRequest struct {
        NfName     string   `json:"nf_name"`
        SbiAddress string   `json:"sbi_address"`
        SbiPort    int      `json:"sbi_port"`
        Services   []string `json:"services"`
    }
    ```

- POST Response

    ```go
    type NnsfNfmPostResponse struct {
        ServiceToken string `json:"service_token"`
        ServiceUuid  string `json:"service_uuid"`
    }
    ```

- DELETE Request

    ```go
    type NnsfNfmDeleteRequest struct {
        NfName   string   `json:"nf_name"`
        Services []string `json:"services"`
    }
    ```

### nnsf-disc

- GET Request

    ```go
    type NnsfDiscGetRequest struct {
        Service string `json:"service"`
    }

    type NnsfDiscGroupGetRequest struct {
        Services []string `json:"services"`
    }
    ```

- GET Response

    ```go
    type NnsfDiscGetResponse struct {
        Service    string `json:"service"`
        NfName     string `json:"nf_name"`
        SbiAddress string `json:"sbi_address"`
        SbiPort    int    `json:"sbi_port"`
    }

    type NnsfDiscGroupGetResponse struct {
        Services []NnsfDiscGetResponse `json:"services"`
    }
    ```

## Implementation and verification notes

See [IMPLEMENTATION.md](IMPLEMENTATION.md) for persistence design, MongoDB configuration, and reproducible supplemental HTTP tests.
The original push and pull-request checks remain enabled; each workflow also supports manual verification from Actions.


### Enable Actions after forking

Open the fork's **Actions** tab and, if GitHub displays the fork-workflow notice,
click **I understand my workflows, go ahead and enable them**. Then a normal
commit pushed to `main` should create the three workflow runs. Check their commit
SHA and event type instead of treating a successful manual run as proof of an
automatic push trigger.
