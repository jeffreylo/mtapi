#!/bin/bash

basedir="${GOPATH}/src/github.com/jeffreylo/mtapi"
mkdir ${basedir}/tmp
curl -sS http://web.mta.info/developers/data/nyct/subway/google_transit.zip > ${basedir}/tmp/mta_gtfs.zip
unzip ${basedir}/tmp/mta_gtfs.zip -d ${basedir}/tmp
rm ${basedir}/tmp/mta_gtfs.zip
cp -a ${basedir}/tmp/. ${basedir}/data/gtfs
rm -rf ${basedir}/tmp
