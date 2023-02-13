package templates

const Configuration = `# global variables
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
`
