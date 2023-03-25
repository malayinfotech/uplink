#!/usr/bin/env bash
set -ueo pipefail
set +x

# This is helper script to setup and run postgres docker image with
# all configuration needed to execute storx-sim integration tests locally
# Script is doing:
# * starts postgres docker container
# * waits for DB to be available
# * creates database for storx-sim

LOG_FILE=${STORX_SIM_POSTGRES_LOG:-"storx-sim-postgres.log"}
CONTAINER_NAME=storx_sim_postgres

cleanup(){
  docker rm -f $CONTAINER_NAME
}
trap cleanup EXIT

docker run --rm -d -p 5433:5432 --name $CONTAINER_NAME -e POSTGRES_PASSWORD=tmppass postgres:12.3 -c log_min_duration_statement=0
docker logs -f $CONTAINER_NAME > $LOG_FILE 2>&1 &

STORX_SIM_DATABASE=${STORX_SIM_DATABASE:-"teststorx"}

RETRIES=10

until docker exec $CONTAINER_NAME psql -h localhost -U postgres -d postgres -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 1
done

docker exec $CONTAINER_NAME psql -h localhost -U postgres -c "create database $STORX_SIM_DATABASE;"

export STORX_SIM_POSTGRES="postgres://postgres:tmppass@localhost:5433/$STORX_SIM_DATABASE?sslmode=disable"