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

