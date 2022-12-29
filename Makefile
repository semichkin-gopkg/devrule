# HelperRules
CloneIfNotExists: 
	[ -d "${to}" ] || git clone ${repo} ${to}

# ServiceRules
AutomationCore_Load: 
	make CloneIfNotExists \
	repo="git@git.collabox.dev:nse/automation/core/local_env.git" \
	to="services/automation/core/local_env"

AutomationCore_Init: 
	cd services/automation/core/local_env && make Init

AutomationCore_ApplyMigrations: 
	cd services/automation/core/local_env && make ApplyMigrations || true

AutomationCore_Start: 
	cd services/automation/core/local_env && make Start

AutomationCore_Stop: 
	cd services/automation/core/local_env && make Stop

AutomationCore_Restart: 
	cd services/automation/core/local_env && (make Restart || make Stop && make Start)

NatsWatcher_Load: 
	make CloneIfNotExists \
	repo="git@git.collabox.dev:nse/nats/watcher/local_env.git" \
	to="services/nats/watcher/local_env"

NatsWatcher_Init: 
	cd services/nats/watcher/local_env && make Init

NatsWatcher_ApplyMigrations: 
	cd services/nats/watcher/local_env && make ApplyMigrations || true

NatsWatcher_Start: 
	cd services/nats/watcher/local_env && make Start

NatsWatcher_Stop: 
	cd services/nats/watcher/local_env && make Stop

NatsWatcher_Restart: 
	cd services/nats/watcher/local_env && (make Restart || make Stop && make Start)

# MainRules
Load: AutomationCore_Load NatsWatcher_Load

Init: AutomationCore_Init NatsWatcher_Init

ApplyMigrations: AutomationCore_ApplyMigrations NatsWatcher_ApplyMigrations

Start: AutomationCore_Start NatsWatcher_Start

Stop: AutomationCore_Stop NatsWatcher_Stop

Restart: AutomationCore_Restart NatsWatcher_Restart

