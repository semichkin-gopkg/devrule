package templates

const Configuration = `# global variables
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
`
