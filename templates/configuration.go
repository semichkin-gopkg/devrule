package templates

const Configuration = `# global variables
GV:
  RepoBase: "https://github.com/semichkin-gopkg"
  LoadingFolder: "services"

GlobalRules:
  Test: >
    echo "test"

MainRules:
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
    Rules:
      Load: "git clone {some_2}"
  # etc...

DefaultServiceRules:
  Load: >
    make _clone \
    repo="{{GV.RepoBase}}/{{V.Path}}.git" \
    to="{{GV.LoadingFolder}}/{{V.Path}}"

  Actualize: >
    make {{V.ServiceName}}_Load &&
    cd {{GV.LoadingFolder}}/{{V.Path}} &&
    git pull origin $(git branch --show-current)
`
