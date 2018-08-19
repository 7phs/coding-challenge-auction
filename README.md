## Requirements

1. Go 1.11 beta3+

2. Libraries which used in a project and downloaded by ```go get`` command:
     
     - [gin](https://github.com/gin-gonic/gin) - micro-services framework
     - [logrus](https://github.com/sirupsen/logrus) - logger frameworks
     - [cobra](https://github.com/spf13/cobra) - tools helping developing CI

## Build

### Manual

A developer should load all dependencies library before building a binary, run command:
```bash
GO111MODULE=on go get
``` 

Then building an executable:
```bash
GO111MODULE=on go build -o ./auction
``` 

### Using makefile

Run a command
```bash
make build
```

### Using docker

Run a command
```bash
make image
```

## Testing

### Manual

A developer should load all dependencies library before testing, run command:
```bash
GO111MODULE=on go get
``` 

Then building an executable:
```bash
GO111MODULE=on go test ./...
``` 

### Using makefile

Run a command
```bash
make testing
```

## Benchmarking

### Manual

A developer should load all dependencies library before testing, run command:
```bash
GO111MODULE=on go get
``` 

Then building an executable:
```bash
GO111MODULE=on go test ./... -bench . -run ^$
``` 

### Using makefile

Run a command
```bash
make bench
```

## Run

### Quick

Run a service using a makefile command
```bash
make run
```

Or using built docker image

```bash
make image-run
```

### Manual

Environment variable:

- **ADDR** - address which a service listening. Default: ':8080'. Example: ':9090;

- **CORS** - switching on response headers supporting CORS. Default: 'false'. Supported: 'true' or 'false'

- **STAGE** - a name of running staging. Default: 'development'. Supported: 'production', 'development', 'testing'

- **LOG_LEVEL** - a level of logging messages. Default: 'info'. Supported: 'debug', 'info', 'warning', 'error'

Command to run:
```bash
ADDR=:8080 STAGE=production LOG_LEVEL=info ./auction run 
```

## Console commands

Environment variable:

- **ADDR** - a service's URL. Default: 'http://localhost:8080'. Example: 'http://auction.dev'

A service supporting a several command to communicate with a service using REST API:

- **bid push** --user_id [USER_ID] --item_id [ITEM_ID] --bid 134.45. 
   
  _Note_: possible push a bid for any user and item id which before not added to service.
  
- **item get** --id [ITEM_ID]

- **item top** --id [ITEM_ID]

- **item add** --title [TITLE]

- **item update** --id [ITEM_ID] --title [TITLE]
  
- **user get** --id [USER_ID]

- **user add** --name [NAME]

- **user update** --id [ITEM_ID] --name [NAME]

## Performance

Estimated just an internal storage without checking a performance of REST API on MacBool Pro 15 Mid 2012:

- **4.05-6.3 μs/op** - pushing a bid into a storage;
- **75-80 ns/op**    - pushing a bid into an item with heap sorting;
- **720-850 ns/op**  - updating a value of a bid

## Data structures

Almost of models are implementing without synchronization.
Just a item and a user information (but not a bids list) were using Go's mutex.
It helps solving a synchronization bottleneck, in my opinion.

### Sync.Map

A lock-free implementation of a map. Using to store a list of users, items and bids.
It was choosing related no requirements to implemented getting a list of users or items.

Sorting bids of item or bid of user by a bid value or by updated time was implementing in handlers
to preventing synchronization a storage and preventing a bottleneck while pushing bids.    

### Heap

Using to get on top the best bid of an item. 

