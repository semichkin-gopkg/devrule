package templates

const Configuration = `# global variables
GV:
  RepoBase: "https://github.com/semichkin-gopkg"
  LoadingFolder: "services"

GlobalRules:
  Start: "cd docker && docker-compose up -d --build"
  Stop: "cd docker && docker-compose down"
  Restart: "make Stop && make Start"
  
  Env: "echo ${DEFAULT_VALUE} && echo ${EXAMPLE}"

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
  - Name: Configurator
    # service variables
    V:
      Path: "configurator"
  - Name: Promise
    V:
      Path: "promise"
    Rules:
      Load: "git clone {some_2}"
  # etc...
`
