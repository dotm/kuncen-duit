This is an end-to-end test for backend only.
This means there'll be no UI automation.
There'll only be HTTP requests to hit backend API.

Benefits:

- Testing backend functionality.
- Documenting API flow (what APIs need to be hit to achieve a user story).
- Manually hitting an API in dev/staging/prod for operation purposes.

The `single` directory is filled with a main function that hit one API endpoint only.
Useful for local testing when developing one API.

The `multiple` directory is filled with a main function that hit multiple API endpoints.
Useful for end-to-end testing scenarios that simulate business flow of user stories.

## E2E Backend Run

- `cd e2e-backend`
- `go mod tidy`
  - If the backend under test is modified, you might need to exec this to remove lint error.
- To run ping: `go run single/ping/*.go`

## E2E Backend Init

- go mod init kuncenduit-e2e-backend
- add `replace kuncenduit-backend => ../backend` to go.mod
- go get kuncenduit-backend