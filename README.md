## Development Makefile builder
A tool for generating rules for managing a large number of local microservices

### Installing
`go install github.com/semichkin-gopkg/devrule/cmd/devrule@v0.0.19`

### Usage
## Build
`devrule build -c path/to/configuration.[yaml|json] -o path/to/output/Makefile`
## Init
`devrule init -o example/configuration.yaml`

#### configuration.yaml
```yaml
# global variables
GV:
  RepoBase: "https://github.com/semichkin-gopkg"
  LoadingFolder: "services"

EnvFiles:
  - import.env

GlobalRules:
  Start: "cd docker && docker-compose up -d --build"
  Stop: "cd docker && docker-compose down"
  Restart: "make Stop && make Start"

  Env: "echo ${FROM_ENV} && echo ${FROM_LOCAL_ENV} && echo ${FROM_IMPORT_ENV}"

MainRules:
  - "Load"
  - "Actualize"

DefaultServiceRules:
  Load: >
    make _clone \
    repo="{{GV.RepoBase}}/{{V.Path}}.git" \
    to="{{GV.LoadingFolder}}/{{V.Path}}" &&
    cd {{GV.LoadingFolder}}/{{V.Path}} &&
    (make Load || true)
  Actualize: >
    make {{V.ServiceName}}_Load &&
    cd {{GV.LoadingFolder}}/{{V.Path}} &&
    git pull origin $(git branch --show-current) &&
    (make Actualize || true)

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
      Load: "git clone {some_2}"
      Unique: "echo 'test'"
  # etc...
```

#### Run Makefile generation
`devrule -c example/configuration.yaml -o example/Makefile`

#### Result
```makefile
ifneq (,$(wildcard .env))
	include .env
	export
endif
ifneq (,$(wildcard .local.env))
	include .local.env
	export
endif
ifneq (,$(wildcard import.env))
	include import.env
	export
endif


# GlobalRules
_clone: 
	[ -d '${to}' ] || git clone ${repo} ${to}

Env: 
	echo ${FROM_ENV} && echo ${FROM_LOCAL_ENV} && echo ${FROM_IMPORT_ENV}

Restart: 
	make Stop && make Start

Start: 
	cd docker && docker-compose up -d --build

Stop: 
	cd docker && docker-compose down


# ServiceRules
Env_Load: 
	make _clone \ repo="https://github.com/semichkin-gopkg/env.git" \ to="services/env" && cd services/env && (make Load || true)

Env_Actualize: 
	make Env_Load && cd services/env && git pull origin $(git branch --show-current) && (make Actualize || true)

Configurator_Load: 
	make _clone \ repo="https://github.com/semichkin-gopkg/configurator.git" \ to="services/configurator" && cd services/configurator && (make Load || true)

Configurator_Actualize: 
	make Configurator_Load && cd services/configurator && git pull origin $(git branch --show-current) && (make Actualize || true)

Promise_Load: 
	git clone {some_2}

Promise_Unique: 
	echo 'test'

Promise_Actualize: 
	make Promise_Load && cd services/promise && git pull origin $(git branch --show-current) && (make Actualize || true)


# GroupedRules

# Main Rules
Load: Env_Load Configurator_Load Promise_Load
Actualize: Env_Actualize Configurator_Actualize Promise_Actualize

# Namespace1 Rules
Namespace1_Load: Env_Load Configurator_Load Promise_Load
Namespace1_Actualize: Env_Actualize Configurator_Actualize Promise_Actualize

# Namespace2 Rules
Namespace2_Load: Env_Load Configurator_Load
Namespace2_Actualize: Env_Actualize Configurator_Actualize
```
