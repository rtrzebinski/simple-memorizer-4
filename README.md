### simple-memorizer-go

WIP - simple-memorizer-go web app built with go-app

```
$ make
target                         help
------                         ----
help                           Show this help
start                          Start containers (docker compose up)
stop                           Stop containers (docker compose down)
restart                        Stop and start containers
destroy                        Stop containers and remove volumes
ps                             Show running containers
migrate                        Run db migrations (migrate up)
migrate-down                   Revert db migrations (migrate down)
migrate-drop                   Drop db without confirmation (migrate drop)
seed                           Seed the database with some example data
reseed                         Destroy, recreate and seed the database (no confirmation)
db                             Db CLI client connection
build                          Build client and server
run                            Build and run locally
test                           Test all
test-short                     Test short (unit)
dev                            Prepare dev environment (stop + start + migrate + seed)
```

### Documentation

https://go-app.dev/getting-started

[Dev documentation - local environment](https://github.com/rtrzebinski/simple-memorizer-go/wiki/Dev-documentation---local-environment) 
