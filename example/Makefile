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


# InternalRules
_git_pull: 
	@[ -d '${to}' ] || git clone ${repo} ${to} && git --git-dir=${to}/.git --work-tree=${to} pull origin $(shell git --git-dir=${to}/.git --work-tree=${to} branch --show-current)


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
Env_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/env.git" to="services/env"

Configurator_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/configurator.git" to="services/configurator"

Promise_Unique: 
	echo 'test'

Promise_Pull: 
	@make -f ${mk} _git_pull  repo="https://github.com/semichkin-gopkg/promise.git" to="services/promise"


# GroupedRules

# Main Rules
Pull: Env_Pull Configurator_Pull Promise_Pull

# Namespace1 Rules
Namespace1_Pull: Env_Pull Configurator_Pull Promise_Pull

# Namespace2 Rules
Namespace2_Pull: Env_Pull Configurator_Pull

