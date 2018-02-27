# mtapi

[![Build Status](https://travis-ci.org/jeffreylo/mtapi.svg?branch=master)](https://travis-ci.org/jeffreylo/mtapi)

Wrapping the MTA's GTFS API, aggregating stops.txt into stations (inspired by
[jonthorton's approach](https://github.com/jonthornton/MTAPI)).

## Dependencies

- go
- dep
- [an MTA API key](http://datamine.mta.info/user)

## Getting Started

```
$ brew bundle
$ open http://datamine.mta.info/user
$ git clone git@github.com:jeffreylo/mtapi
$ cd $GOPATH/src/github.com/jeffreylo/mtapi
$ go install ./...
$ mtapi -apiKey=$MTA_API_TOKEN -path=$(pwd)/data/gtfs -port=8080
$ echo '{"jsonrpc": "2.0","method": "GetStation","params": { "ID": "132" },"id": "243a718a-2ebb-4e32-8cc8-210c39e8a14b"}' | http POST http://localhost:8080/rpc
```
