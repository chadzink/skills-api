## Code Path

```mermaid
graph
  USR{API User} --> WAPI(Web API)
  subgraph API
    WAPI <--> ROUTER[Router -routes.go]
    ROUTER <--> HANDLER[Handler -handler pkg]
    HANDLER <--> DAL[Data Access Layer -database pkg]
    DAL <--> MODELS[Models -models pkg]
  end
```

### Example code path
Request Skill by Id `/skill/1`: `main.go` --> `router.go` --> `handler/skills.go` --> `database/skills.go` --> `models/skill.go`

## Authentication
TO DO: Add Authentication diagram for create user, create API key, & login.
