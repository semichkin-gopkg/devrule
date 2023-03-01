# Expressions
mk := $(abspath $(lastword $(MAKEFILE_LIST)))
mkdir := $(dir $(mk))
pwd := $(shell pwd)
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


# GlobalRules
Echo: 
	@echo ${message}

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
Env_Info: 
	@make -f ${mk} Echo message="Name: Env"

Configurator_Info: 
	@make -f ${mk} Echo message="Name: Configurator"

Promise_Unique: 
	echo 'test'

Promise_Info: 
	@make -f ${mk} Echo message="Name: Promise"


# GroupedRules

# Main Rules
Info: Env_Info Configurator_Info Promise_Info

# Namespace1 Rules
Namespace1_Info: Env_Info Configurator_Info Promise_Info

# Namespace2 Rules
Namespace2_Info: Env_Info Configurator_Info

