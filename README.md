# Git API Service# Description

This service fetches data from Git APIs to retrieve repository commits, saves the data in a persistent store, and continuously monitors the repositories for changes at a set interval.

## Requirements- Golang 1.20+
- PostgreSQL

## 1. Clone the repository, cd into the project folder and download required go dependencies
```bash
git clone https://github.com/kenmobility/git-api-service.git
```
```bash
cd git-api-service
```
```bash
go mod tidy
```
## 2. Run below command to Duplicate .env.example file and rename to .env
```bash
cp .env.example .env
```

## 3. Set environmental variables:
- the duplicated .env.example file already has default variables that the program needs to run except for GIT_HUB_TOKEN env variable.
- The program can run without GIT_HUB_TOKEN variable, but with a rate limit of just 60 requests within a time frame, but to extend the rate limit to 5000 requests, a valid GitHub token should be added to the .env file. 
- (OPTIONAL) Go to [https://github.com/](GitHub) to set up a GitHub API token (i.e Personal access token) and set the value for the GIT_HUB_TOKEN environmental variable on the .env file.

- (OPTIONAL) if the default values of the DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD or DATABASE_NAME in .env file were altered, ensure that the 'make postgres' command in the makefile matches the new set values.

## 4 Open Docker desktop application
- Ensure that docker desktop is started and running on your machine 

## 5. Run makefile commands 
- run 'make postgres' to pull and run PostgreSQL instance as docker container
```bash
make postgres
```
- run 'make createdb' to create a database
```bash
make createdb
```

## 6. Unit Testing
Run 'make test' to run the unit tests:
```bash
make test
```

## 7. Start web server
- run 'make server' to start the service
```bash
make server
```

## 8. Endpoint requests
- POST application/json Request to add a new repository
``` 
curl -d '{"name": "GoogleChrome/chromium-dashboard"}'\
  -H "Content-Type: application/json" \
  -X POST http://127.0.0.1:5000/repository \
```

- GET Request to fetch all the repositories on the database
```
curl -L \
  -X GET http://127.0.0.1:5000/repositories \
```

- GET Request to fetch all the commits fetched from github API for any repo using repository Id 
```
curl \
  -X GET http://127.0.0.1:5000/repos/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a/commits \
```

- GET Request to get repository metadata using repository id. 
``` 
curl -L \
  -X GET http://127.0.0.1:5000/repository/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a \
```

- GET Request to fetch N (as limit) top commit authors of the any added repository using its repository id with limit as query param, if limit is not passed, a defualt limit of 10 is used.
```
curl -L \
  -X GET http://127.0.0.1:5000/repos/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a/top-authors?limit=5 \
```

## Clean Slate: Removing Database
- To remove all the data (droping the database) run 'make dropdb' to delete the database, ensure program isn't running and database isn't openned in any database client e.g pgAdmin or TablePlus, etc
```bash
make dropdb
```
- Then run 'make createdb' to recreate the database before restarting the program with 'make server'
```bash
make createdb
```
- Then run 'make server' to rerun database migrations (recreate all tables, seeds default repo) and starts program
```bash
make server
```