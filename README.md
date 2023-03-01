## Development Makefile builder
A tool for generating rules for managing a large number of local microservices

### Installing
`go install github.com/semichkin-gopkg/devrule/cmd/devrule@v0.0.23`

### Usage
## Build
`devrule build -c path/to/configuration.[yaml|json] -o path/to/output/Makefile`
## Init
`devrule init -o example/configuration.yaml`

#### configuration.yaml
```yaml
# global variables
GV:
  Repo: "https://github.com/semichkin-gopkg"
  ServiceDir: "services"

Expressions:
  - export FROM_EXPRESSIONS := $(shell echo "from_expressions")

EnvFiles:
  - import.env

GlobalRules:
  Echo: "@echo ${message}"
  Start: "cd docker && docker-compose up -d --build"
  Stop: "cd docker && docker-compose down"
  Restart: "make Stop && make Start"
  EnvFilesTest: "@echo ${FROM_IMPORT_ENV}"
  ExpressionsTest: "@echo ${FROM_EXPRESSIONS}"

MainRules:
  - "Info"

DefaultServiceRules:
  Info: >
    @make -f ${mk} Echo
    message="Name: {{V.ServiceName}}"

Services:
  - Name: Env
    Groups: ["_all"] # group _all tells that service rules should be included to all other groups
    # service variables
    V:
      Path: "env"
  - Name: Configurator
    Groups: ["Namespace1", "Namespace2"]
    V:
      Path: "configurator"
  - Name: Promise
    Groups: ["Namespace1"]
    V:
      Path: "promise"
    Rules:
      Unique: "echo 'test'"
  # etc...
```

#### Run Makefile generation
`devrule -c example/configuration.yaml -o example/Makefile`

#### Result
```makefile
# Expressions
mk := $(abspath $(lastword $(MAKEFILE_LIST)))
mkdir := $(dir $(mk))
pwd := $(shell pwd)
export FROM_EXPRESSIONS := $(shell echo "from_expressions")


# EnvFiles
ifneq (,$(wildcard .env))
	include .env
	export
endif
ifneq (,$(wildcard import.env))
	include import.env
	export
endif


# GlobalRules
Echo: 
	@echo ${message}

EnvFilesTest: 
	@echo ${FROM_IMPORT_ENV}

ExpressionsTest: 
	@echo ${FROM_EXPRESSIONS}

Restart: 
	make Stop && make Start

Start: 
	cd docker && docker-compose up -d --build

Stop: 
	cd docker && docker-compose down


# ServiceRules
Env_Info: 
	@make -f ${mk} Echo message="Name: Env"

Configurator_Info: 
	@make -f ${mk} Echo message="Name: Configurator"

Promise_Unique: 
	echo 'test'

Promise_Info: 
	@make -f ${mk} Echo message="Name: Promise"


# GroupedRules

# Main Rules
Info: Env_Info Configurator_Info Promise_Info

# Namespace1 Rules
Namespace1_Info: Env_Info Configurator_Info Promise_Info

# Namespace2 Rules
Namespace2_Info: Env_Info Configurator_Info
```