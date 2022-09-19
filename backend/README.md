## Backend Todo

- continue user.me
- Use sqlc? Use sqlx?
- add event sourcing especially to user.signup
- Catch panic or error.
- Continue implement user.me, user.signout
- Connect with frontend
- Deploy to Lambda
- Readd nonce for optimistic locking?
- Use unique on stream id and nonce to avoid racing
- Retry whole function on race
- Differentiate between individual user and company user?
- Encapsulate db acquire connection and return close function.
- Use zerolog Errs instead of Err to log errors (?)
- Add optional send to sns in logger instance.
- If not local, send to dynamodb and sns instead of standard error.

## Run Backend in Local

- go run *.go

## Test Backend in Local

## Deploy Backend

While local deployment use `go run *.go` to init the http server,
deployment in other stages should be from Terraform into AWS Lambda and API Gateway.

## Migrate Database

- Install golang-migrate/migrate using Linux (*.deb package)
  - https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#linux-deb-package
    - curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
    - echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
    - apt-get update
    - apt-get install -y migrate
- Create migration files:
  - `migrate create -ext sql -dir ./migration -seq create_users`
- Migrating:
  - Export database URL: `export DATABASE_URL="postgres:Test123@$WIN_HOST_IP:5432/postgres?sslmode=disable"`
    - postgres://username:password@localhost:5432/database_name
  - Up to latest: `migrate -database ${DATABASE_URL} -path ./migration up`
  - Up 1: `migrate -database ${DATABASE_URL} -path ./migration up 1`
  - Down all: `migrate -database ${DATABASE_URL} -path ./migration down`
  - Down 1: `migrate -database ${DATABASE_URL} -path ./migration down 1`
- for local, add `?sslmode=disable` in DATABASE_URL
  - this means that the connection with our database will not be encrypted

## Backend Init

- cd backend
- go mod init kuncenduit-backend

## Local PostgreSQL

For local testing, install postgresql:

- sudo apt update && sudo apt upgrade
- sudo apt install postgresql postgresql-contrib
  - the available version in the default package repository might not be up to date.
  - search the internet to find out how to download the latest version.
  - for version 14 on Ubuntu 20: https://techviewleo.com/how-to-install-postgresql-database-on-ubuntu/
- psql --version
- sudo passwd postgres
  - set password for postgres database

Using local postgresql with psql:

- sudo service postgresql status
- sudo service postgresql start
  - sudo service postgresql stop
- sudo -u postgres psql
- common psql meta-commands:
  - `\l` to list databases
  - `\c database_name` to connect to a database
  - `\dt` to list tables in current database
  - `\q` or `exit`

## Finding Local WSL IP

- get host ip from wsl:
  - ipconfig.exe | grep 'vEthernet (WSL)' -A4 | cut -d":" -f 2 | tail -n1 | sed -e 's/\s*//g'
- get wsl ip from host:
  - wsl hostname -i

## Project Structure

- http is filled with any HTTP request related stuff: API endpoint path, JSON body, handler, etc.

## Event Sourcing Read List

https://www.eventstore.com/event-sourcing
https://www.eventstore.com/blog/event-immutability-and-dealing-with-change
https://www.eventstore.com/blog/event-sourcing-and-cqrs

## Other Read List

https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda

## Strive for Zero Down Time Deployment

- Be careful when migrating schema since it'll acquire a full table access lock transaction.
- Use versions in codebase, environment variables, etc. to enable blue green or multi version deployment.
