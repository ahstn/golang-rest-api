# Go Rest API
[![Go Report Card](https://goreportcard.com/badge/github.com/phazyy/golang-rest-api)](https://goreportcard.com/report/github.com/phazyy/golang-rest-api)

Example Rest API in Go using [Gin], [SqlBoiler], [Log15] and [Postgres]


## Todo
- [x] Basic CRUD Functionality
- [ ] Get entity relationship data
- [ ] Validation of input (409 if entity exists, 400 if invalid req)
- [ ] API Auth
- [ ] HTTP error handling (returning 404, 403, Adding status code to response payload)
- [ ] Testing and code coverage
- [ ] Handling encrypted fields (e.g. Passwords)
- [ ] Encrypted payloads/json ?
- [ ] Containerise binary
- [x] Review logging (unify output if possible)

[Gin]: https://github.com/gin-gonic/gin
[SqlBoiler]: https://github.com/vattle/sqlboiler
[Log15]: https://github.com/inconshreveable/log15
[Postgres]: https://github.com/postgres/postgres
