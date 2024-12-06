#!/usr/bin/env bash

CLUSTER_NAME='kafka-dev'

REPOSITORY='mihkels'
APP_NAME='kafka-tester'
VERSION='2024.11.11.62'
APP_TYPE=('consumer' 'producer')

for type in "${APP_TYPE[@]}"; do
    kind load docker-image --name ${CLUSTER_NAME} ${REPOSITORY}/${APP_NAME}:"${type}"-${VERSION}-golang
done