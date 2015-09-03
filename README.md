#WIP

### Neutrino is under heavy development, there are still some rough edges.

# REALBASE
[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

[![Build Status](https://travis-ci.org/go-realbase/realbase.svg?branch=master)](https://travis-ci.org/go-realbase/realbase)
--------------
Realbase is hybrid realtime-REST backend written in Go with MongoDB as a database.

### Any kind of feedback or contribution is welcomed warmly.

# Docs

Interactive documentation can be found [here](http://docs.realbas3.apiary.io/#reference).

# Development
--------------

```bash
$ go get github.com/go-neutrino/neutrino-core
```

In the root of the project - `$GOPATH/src/github.com/go-neutrino/neutrino` you can execute the following:

```bash
$ go run neutrino.go
```

```bash
$ go test -v ./..
```

# Goals for the initial release

- [x] MongoDB to store the data
- [x] User management with JWT 
- [ ] Landing page (http://neutrinoapp.com)
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
  - [x] Integration tests
- [ ] WebSockets API - Realtime API
  - [ ] Types support
  - [ ] Javascript SDK
