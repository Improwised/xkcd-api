# Xkcd+PoorlyDrawnLines

### **Index**
- [Run Locally](#run-locally)

- [Run in Docker](#run-in-docker)

- [Kick Start Commands](#kick-start-commands)

---
### **Run Locally**
- Pre-requictic: golang should be installed.
- Clone this Repository. 
- Open terminal and go inside this Repository
- To sync go packages you need to run following command.

    ### `go mod vendor`
- Copy `.env.example` as `.env` for environment variables.
- Following command is use to remove unused package from your `go.mod` and `go.sum` files.(I fchange any package)

    ### `go mod tidy`
- **make start** : To start api, it basically runs `go run app.go`

---
## Run in Docker
## Execution

1. Run ```docker-compose up```
