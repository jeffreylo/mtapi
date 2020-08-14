#!/bin/bash
set -eux

basedir="$( cd "$( dirname "$0" )/.." && pwd )"
rm -rf ${basedir}/tmp
mkdir -p ${basedir}/tmp
curl -sS http://web.mta.info/developers/data/nyct/subway/google_transit.zip > ${basedir}/tmp/mta_gtfs.zip
unzip ${basedir}/tmp/mta_gtfs.zip -d ${basedir}/tmp
rm ${basedir}/tmp/mta_gtfs.zip
mkdir -p ${basedir}/mta/testdata/gtfs
cp -a ${basedir}/tmp/. ${basedir}/mta/testdata/gtfs
