#WIP

### Neutrino is under heavy development, there are still some rough edges.

# NEUTRINO

![Neutrino](https://media.giphy.com/media/3o85xnGaP3m49VmBDW/giphy.gif)


[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

[![Join the chat at https://gitter.im/go-neutrino/neutrino-core](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-neutrino/neutrino-core?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)   
[![Build Status](https://travis-ci.org/go-neutrino/neutrino-core.svg?branch=master)](https://travis-ci.org/go-neutrino/neutrino-core)
--------------
Neutrino is hybrid realtime-REST backend written in Go with MongoDB as a database.

### Any kind of feedback or contribution is welcomed warmly.

# Docs

Interactive documentation can be found [here](http://docs.realbas3.apiary.io/#reference).

# Development
--------------

### You need NATS and MongoDB running

```bash
$ go get github.com/go-neutrino/neutrino
```

In the root of the project - `$GOPATH/src/github.com/go-neutrino/neutrino` you can execute the following:

```bash
$ go run api-service/main.go
$ go run realtime-service/main.go
$ go run queue-broker-service/main.go
```

```bash
$ go test -v ./..
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
- [] WebSockets API - Realtime API
  - [ ] Types support
  - [ ] Javascript SDK

# TODO
#### Updated on the go:

- [ ] Api service should send realtime-jobs to the queue for the broker and realtime service to process them
- [ ] Cover api service's new functionality, realtime service and queue-broker with tests
- [ ] Maybe abstract the websockets logic to a separate plugin
