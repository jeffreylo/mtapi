# mtapi

[![Build
Status](https://travis-ci.org/jeffreylo/mtapi.svg?branch=master)](https://travis-ci.org/jeffreylo/mtapi)

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
$ mtapi -api-key=${MTA_API_TOKEN} -gtfs-path=$(pwd)/data/gtfs -port=9090
$ open http://localhost:9090
```

## Demo

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)
