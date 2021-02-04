# Real-time WebSocket sport events dashboard

The entire solution is packed into Docker, docker-compose.yaml starts up the application, docker-compose.test.yaml runs tests only, separate `Dockerfile` and `Dockerfile.test` & `Docker-compose.yaml` and `Docker-compose.test.yaml` were used for this purpose, they use different databases. Dockerized app contains PostgreSQL, pgAdmin, Golang backend + demo client, React application. 

# Starting the project

Execute the following commands under `/app` directory:

#### - Run the application: `docker-compose up --build`
This will install and migrate the database, server, demo client, and the React application.
React application is accessible at http://localhost:3000/

#### - Run tests: `docker-compose -f .\docker-compose.test.yaml up --build`
This will run all the test against the test database inside transactions.

# Important setup step

Self-signed certificates are used in this solution, browser will block frontend request to backend, in order to add a cert into browser exceptions:
1) Open web-browser
2) Go to https://localhost:8000/
3) Click "Accept risk and continue" that will add the certificate into exceptions

# To-do things
Cached results flushing (out of scope for now).
* Remove old results from the frontend state
* Remove old results from the Go backend the latest results slice

The Simplest implementation is to delete/overwrite older results indexes when there are e.g. more than 1000 values, but the real load is not clear yet, so these are my assumptions.

# Structure and approach
## Backend

CQRS - Command Query Responsibility Segregation, commands and queries are split into separate files.

DDD (Domain Driven Development) - common Go directory with the shared domain details, these are `commands`, `queries`, `models`, `errors`, there are also common domain errors and exclusive domain model errors.
The "stateful" look like models are used to mitigate type errors and make code look more like a written documentation, this reduces the layer of abstraction, that's what the DDD is about.

BDD (Behavior-driven Development) - tests were done the BDD style with Ginkgo/Gomega.
Same DB connection for tests asynchronous might cause data interference, for that the Dashboard WebSocket module tests are running inside the transactions. For that reason the WebSocket module tests initiate different DB connections so that they will use different transactions as well and data interference will be impossible.

Monorepo - separation of concerns, backend and frontend code is located in different directories, while the docker-compose allows to make one-liner deployment.

CMD pattern - It's not an official standard defined by the core Go dev team, helps to manage multiple main.go entry-points in the future.

YAML config - the most critical configuration details must be stored statically to mitigate side effects, so that no global ENV could overwrite it in a sudden.

Custom logger - fast and structured logging. 

crt.crt and key.key are required to run the server in TLS, moving API under the https/wss secure connection.

GORM - ORM used in Go project for database migration easiness. 

Protobuf - protobuf for events, so that it would be easy to integrate any RPC/message-queue services.

Context for the server graceful shutdown and channels to write errors from goroutines back to the main method and handle them there.

Transactions - tests are running in transactions and rollback is performed after, so that the db won't get polluted with test data.

Domain errors - custom error types are providing much more information regarding the states, and helps re-using the codebase for error handling.
## Frontend

As the dashboard is the only component used in app the state management was implemented with built-in React state, effects and hooks.
