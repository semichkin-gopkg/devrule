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
Configurator_Load: 
	make _clone \ repo="https://github.com/semichkin-gopkg/configurator.git" \ to="services/configurator" && cd services/configurator && (make Load || true)

Configurator_Actualize: 
	make Configurator_Load && cd services/configurator && git pull origin $(git branch --show-current) && (make Actualize || true)

Promise_Load: 
	git clone {some_2}

Promise_Actualize: 
	make Promise_Load && cd services/promise && git pull origin $(git branch --show-current) && (make Actualize || true)


# GroupedRules

# Main Rules
Load: Configurator_Load Promise_Load
Actualize: Configurator_Actualize Promise_Actualize

# Namespace1 Rules
Namespace1_Load: Configurator_Load Promise_Load
Namespace1_Actualize: Configurator_Actualize Promise_Actualize

# Namespace2 Rules
Namespace2_Load: Configurator_Load
Namespace2_Actualize: Configurator_Actualize

