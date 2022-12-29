## Development Makefile builder
A tool for generating rules for managing a large number of local microservices

### Installing
`go install github.com/semichkin-gopkg/devrule/cmd/devrule@latest`

### Usage
`devrule -c path/to/configuration.[yaml|json] -o path/to/output/Makefile`

### Example
#### configuration.yaml
```yaml
# global variables
GV:
  RepoBase: "https://github.com/semichkin-gopkg"
  LoadingFolder: "services"

HelperRules:
  CloneIfNotExists: >
    [ -d "${to}" ] || git clone ${repo} ${to}

EnabledRules:
  - "Load"
  - "Actualize"
  
Services:
  Configurator:
    # service variables
    V:
      Path: "configurator"
  Promise:
    V:
      Path: "promise"
  # etc...

DefaultServiceRules:
  Load: >
    make CloneIfNotExists \

    repo="{{GV.RepoBase}}/{{V.Path}}.git" \

    to="{{GV.LoadingFolder}}/{{V.Path}}"

  Actualize: >
    make Load && 
    cd {{GV.LoadingFolder}}/{{V.Path}} && 
    git pull origin $(git branch --show-current)
```

#### Run Makefile generation
`devrule -c example/configuration.yaml -o example/Makefile`

#### Result
```makefile
# HelperRules
CloneIfNotExists: 
	[ -d "${to}" ] || git clone ${repo} ${to}

# ServiceRules
Configurator_Load: 
	make CloneIfNotExists \
	repo="https://github.com/semichkin-gopkg/configurator.git" \
	to="services/configurator"

Configurator_Actualize: 
	make Configurator_Load && cd services/configurator && git pull origin $(git branch --show-current)

Promise_Load: 
	make CloneIfNotExists \
	repo="https://github.com/semichkin-gopkg/promise.git" \
	to="services/promise"

Promise_Actualize: 
	make Promise_Load && cd services/promise && git pull origin $(git branch --show-current)

# MainRules
Load: Configurator_Load Promise_Load

Actualize: Configurator_Actualize Promise_Actualize
```
