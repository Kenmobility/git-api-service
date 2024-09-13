# Git API Service 
Description:

This service fetches data from GitHub APIs to retrieve GitHub repository commits, saves the data in a persistent store (PostgreSQL), and continuously monitors and reconciles the fetched commits of all the added repositories at a set interval.

## NOTE: environmental variables:
- the .env.example file already has default variables that the program needs to run except for GIT_HUB_TOKEN env variable.
- The program can run without GIT_HUB_TOKEN variable, but with a rate limit of just 60 requests within a time frame, to extend the rate limit to 5000 requests, a valid GitHub token should be added to the .env file. 
- Go to [https://github.com/](GitHub) to set up a GitHub API token (i.e Personal access token) and set the value for the GIT_HUB_TOKEN environmental variable on the .env file.

## Requirements
- Docker Desktop app

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

## 2. Unit Testing
Run 'make test' to run the unit tests:
```bash
make test
```

## 3 Open Docker desktop application
- Ensure that docker desktop is started and running on your machine 

## 4. Run the application
- run 'make' to run application
```bash
make
```

## 5. Endpoint requests
- POST application/json Request to add a new repository
``` 
curl -d '{"name": "GoogleChrome/chromium-dashboard"}'\
  -H "Content-Type: application/json" \
  -X POST http://localhost:8080/repository \
```

- GET Request to fetch all the repositories on the database
```
curl -L \
  -X GET http://localhost:8080/repositories \
```

- GET Request to fetch all the commits fetched from github API for any repository using repository Id, response is paginated, pass 'limit' and 'page' as query params to get next pages.
```
curl \
  -X GET http://localhost:8080/repos/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a/commits?limit=20&page=1 \
```

- GET Request to get repository metadata using repository id. 
``` 
curl -L \
  -X GET http://localhost:8080/repository/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a \
```

- GET Request to fetch N (as limit) top commit authors of the any added repository using its repository id with limit as query param, if limit is not passed, a defualt limit of 10 is used.
```
curl -L \
  -X GET http://localhost:8080/repos/5846c0f0-81f5-45e3-9d4a-cfc6fe4f176a/top-authors?limit=5 \
```

## Clean Slate: 
Removing containers
- To remove the containers run 'make down'
```bash
make down
```
- To remove pulled images run 'make clean'
```bash
make clean
```
- Then run 'make' to re-pull images, re-start all containers and run the program
```bash
make
```