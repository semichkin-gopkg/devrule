## Development Makefile builder
A tool for generating rules for managing a large number of local microservices

### Installing
`go install github.com/semichkin-gopkg/devrule/cmd/devrule@v0.0.6`

### initializing
`devrule init -o path/to/output/configuration.yaml`

### Usage
`devrule build -c path/to/configuration.[yaml|json] -o path/to/output/Makefile`

### Example

#### Run configuration.yaml initialization
`devrule init -o example/configuration.yaml`

#### configuration.yaml
```yaml
# global variables
GV:
  RepoBase: "https://github.com/semichkin-gopkg"
  LoadingFolder: "services"

GlobalRules:
  Build: "cd docker && docker-compose build"
  Start: "cd docker && docker-compose up -d"
  Stop: "cd docker && docker-compose down"
  Restart: "cd docker && docker-compose restart"

MainRules:
  - "Load"
  - "Actualize"

DefaultServiceRules:
  Load: >
    make _clone \
    repo="{{GV.RepoBase}}/{{V.Path}}.git" \
    to="{{GV.LoadingFolder}}/{{V.Path}}" &&
    cd {{GV.LoadingFolder}}/{{V.Path}} &&
    make Load || true
  Actualize: >
    make {{V.ServiceName}}_Load &&
    cd {{GV.LoadingFolder}}/{{V.Path}} &&
    git pull origin $(git branch --show-current) &&
    make Actualize || true

Services:
  Configurator:
    # service variables
    V:
      Path: "configurator"
  Promise:
    V:
      Path: "promise"
    Rules:
      Load: "git clone {some_2}"
  # etc...
```

#### Run Makefile generation
`devrule -c example/configuration.yaml -o example/Makefile`

#### Result
```makefile
# GlobalRules
_clone: 
	[ -d '${to}' ] || git clone ${repo} ${to}

Build: 
	cd docker && docker-compose build

Restart: 
	cd docker && docker-compose restart

Start: 
	cd docker && docker-compose up -d

Stop: 
	cd docker && docker-compose down

# ServiceRules
Configurator_Load: 
	make _clone \ repo="https://github.com/semichkin-gopkg/configurator.git" \ to="services/configurator" && cd services/configurator && make Load || true

Configurator_Actualize: 
	make Configurator_Load && cd services/configurator && git pull origin $(git branch --show-current) && make Actualize || true

Promise_Load: 
	git clone {some_2}

Promise_Actualize: 
	make Promise_Load && cd services/promise && git pull origin $(git branch --show-current) && make Actualize || true

# MainRules
Load: Configurator_Load Promise_Load

Actualize: Configurator_Actualize Promise_Actualize
```
