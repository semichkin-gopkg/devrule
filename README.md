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
  Start: "cd docker && docker-compose up -d --build"
  Stop: "cd docker && docker-compose down"
  Restart: "make Stop && make Start"
  EnvFilesTest: "@echo ${FROM_IMPORT_ENV}"
  ExpressionsTest: "@echo ${FROM_EXPRESSIONS}"

MainRules:
  - "Pull"

DefaultServiceRules:
  Pull: >
    @make -f ${mk} _git_pull 
    repo="{{GV.Repo}}/{{V.Path}}.git"
    to="{{GV.ServiceDir}}/{{V.Path}}"

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


# InternalRules
_git_pull: 
	@[ -d '${to}' ] || git clone ${repo} ${to} && git --git-dir=${to}/.git --work-tree=${to} pull origin $(shell git --git-dir=${to}/.git --work-tree=${to} branch --show-current &> /dev/null)


# GlobalRules
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
Env_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/env.git" to="services/env"

Configurator_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/configurator.git" to="services/configurator"

Promise_Unique: 
	echo 'test'

Promise_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/promise.git" to="services/promise"


# GroupedRules

# Main Rules
Pull: Env_Pull Configurator_Pull Promise_Pull

# Namespace1 Rules
Namespace1_Pull: Env_Pull Configurator_Pull Promise_Pull

# Namespace2 Rules
Namespace2_Pull: Env_Pull Configurator_Pull
```