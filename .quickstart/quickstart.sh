#!/bin/bash

#build
make -C ../ build

# # consul instance
# # -v /Users/<your_name>/tmp/consul:/consul/data \
docker run -d --name dev-consul -p 8500:8500 -p 8600:8600/udp \
-e CONSUL_BIND_INTERFACE=eth0 \
consul agent --bootstrap -server -ui -client=0.0.0.0

# postgres instance
# -v $PWD/pgmountlayouttmp:/var/lib/postgresql/data \
docker run -d -p 5438:5432 \
--name volumed-postgres-layout \
-e POSTGRES_PASSWORD=password \
-e PGDATA=/var/lib/postgresql/data/pgdata \
postgres:12

# run migration
docker build --no-cache --network=host ../migration/goose

# Set consul configuration
root="go-micro-layouts"
appname="animal-store"
# Consul url
url=http://host.docker.internal:8500/v1/kv/${root}
token=$1
requestbody='{
  "server": {
    "name": "animal-store",
    "addr": ":8080"
  },
  "postgres_standby": {
    "addr": "host.docker.internal:5438",
    "dbname": "postgres",
    "login": "postgres",
    "passw": "password",
    "conn_max": 80,
    "conn_lifetime": 10,
    "conn_maxidletime": 10
  },
  "postgres_primary": {
    "addr": "host.docker.internal:5438",
    "dbname": "postgres",
    "login": "postgres",
    "passw": "password",
    "conn_max": 40,
    "conn_lifetime": 300,
    "conn_maxidletime": 0
  },
  "metric": {
    "addr": ":8080"
  }
}'

# Configuration
# --- DO NOT EDIT BELOW ---
setConsulConfig () {
  echo "### Setting ${root}/${appname} as:"
  echo "${requestbody}"
  if [[ "$(curl -sX PUT -H "X-Consul-Token: ${token}" -d "${requestbody}" ${url}/${appname})" == "true" ]]; then
    echo "### ${url}/${appname} is set"
  else
    echo "### ERROR: Cannot set ${url}/${appname}"
    exit 1
  fi
}
setConsulConfig

#run service
../bin/app
