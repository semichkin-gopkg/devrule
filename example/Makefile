# Expressions
export PWD := $(shell pwd)
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
_clone: 
	@[ -d '${to}' ] || git clone ${repo} ${to}


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

