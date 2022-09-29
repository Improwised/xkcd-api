# Get first 10 from Xkcd+ 10 from PoorlyDrawnLines

### **Index**
- [Run Locally](#run-locally)

- [Run in Docker](#run-in-docker)
---
### **Run Locally**
- Pre-requictic: golang version 1. 19 should be installed.
1. Clone this Repository. 
2. Open terminal and go inside this Repository
3. To sync go packages you need to run following command.

    ### `go mod vendor`
4. Copy `.env.example` as `.env` for environment variables.

>**_NOTE:_**  The note content.
  Following command is use to remove unused package from your `go.mod` and `go.sum` files.(If change any package)`go mod tidy`

5  **make start** : To start api, it basically runs `go run app.go api`


---
## Run in Docker
## Execution

1. Run ```docker-compose up```
---
## Api End Points
>**_NOTE:_**  Currently have one single api

['http://127.0.0.1:3000/api/v1/getdata']('http://127.0.0.1:3000/api/v1/getdata')

```shell
 curl --location --request GET 'http://127.0.0.1:3000/api/v1/getdata'
 ```
