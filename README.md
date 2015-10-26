#WIP

### Neutrino is under heavy development, there is a lot more work to be done.

# NEUTRINO

![Neutrino](https://media.giphy.com/media/3o85xnGaP3m49VmBDW/giphy.gif)


[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

[![Join the chat at https://gitter.im/go-neutrino/neutrino-core](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-neutrino/neutrino-core?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)   
[![Build Status](https://travis-ci.org/go-neutrino/neutrino.svg?branch=master)](https://travis-ci.org/go-neutrino/neutrino)
--------------
Neutrino is hybrid realtime-REST backend written in Go with MongoDB as a database.

### Any kind of feedback or contribution is welcomed warmly.

# Docs

Interactive documentation can be found [here](http://docs.realbas3.apiary.io/#reference).

# Development
--------------

### You need NATS and MongoDB running, refer to **dev-setup.sh**

```bash
$ go get github.com/go-neutrino/neutrino
```

In the root of the project - `$GOPATH/src/github.com/go-neutrino/neutrino` you need to run all the services:

```bash
$ go run api-service/main.go
$ go run realtime-service/main.go
$ go run queue-broker-service/main.go
```

To run the unit tests execute:

```bash
$ make test
```

To run the integration tests run:

```bash
#make sure that you have all the services running including nats and mongodb
$ go test ./integration-tests/ 
```

# Goals for the initial release

- [x] MongoDB to store the data
- [x] User management with JWT 
- [x] Landing page (http://neutrinoapp.com)
- [ ] App portal (http://app.neutrinoapp.com)
- [ ] API server (http://api.neutrinoapp.com)
- [x] CI
- [x] Docs
- [x] Postman Collection
- [x] REST API
  - [x] Create Applications
  - [x] Read Applications
  - [x] Delete Applications
  - [x] Update Applications
  - [x] Create types
  - [x] Delete types
  - [x] Insert types data
  - [x] Read types data
  - [x] Update types data
  - [x] Delete types data
  - [x] In-app user management
- [ ] WebSockets API - Realtime API
  - [x] Types support
  - [ ] Javascript SDK

# TODO
#### Updated on the go:

- [ ] Transport token from client to realtime service when communicating in realtime
- [ ] Handle services authentication from/to api/realtime/clients
- [ ] Sort out permissions
- [ ] Js client

User logged in -> Get token from api service and reauthenticate with the realtime service
    - HTTP GET -> /v1/api/token?username=username&password=password
    - Conn.Token = token
    
On authentication-required request
    - Validate token
    - Execute