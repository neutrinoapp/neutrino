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
$ go get github.com/neutrinoapp/neutrino
```

In the root of the project - `$GOPATH/src/github.com/neutrinoapp/neutrino` you need to run all the services:

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
$ make integration
```
