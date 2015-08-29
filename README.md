# REALBASE
--------------
Realbase is hybrid realtime-REST backend written in Go with MongoDB as a database.

# Docs

Interactive documentation can be found [here](http://docs.realbas3.apiary.io/#reference).

# Development
--------------

```bash
$ go get github.com/go-realbase/realbase
```

In the root of the project - `$GOPATH/src/github.com/go-realbase/realbase` you can execute the following:

```bash
$ go run realbase.go
```

```bash
$ go test -v ./..
```

# Goals for the initial release

- [x] MongoDB to store the data
- [x] User management with JWT 
- [ ] Landing page
- [ ] App portal (http://app.realbase.co)
- [ ] CI
- [ ] REST API
  - [x] Create Applications
  - [x] Read Applications
  - [ ] Delete Applications
  - [x] Create types
  - [ ] Edit types
  - [ ] Delete types
  - [x] Insert types data
  - [x] Read types data
  - [ ] Edit types data
  - [ ] Delete types data
  - [ ] In-app user management
  - [x] Integration tests
- [ ] WebSockets API - Realtime API
  - [ ] Types support
  - [ ] Javascript SDK
