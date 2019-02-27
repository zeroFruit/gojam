#!/bin/bash

HOST=$1
RATE=$2
DURATION=$3
OUTPUT=$4

USAGE="usage: ./benchmark <HOST> <RATE> <DURATION> <OUTPUT>"

function validateParameters {
    if [ "$1" == "" ] || [ "$2" == "" ] || [ "$3" == "" ]; [ "$4" == "" ];then
        echo ${USAGE}
        exit 1
    fi
}

validateParameters ${HOST} ${RATE} ${DURATION} ${OUTPUT}

go run cmd/traffic/traffic.go POST ${HOST} | \
vegeta attack -rate=${RATE} -duration=${DURATION}s | \
tee ${OUTPUT} | \
vegeta report
