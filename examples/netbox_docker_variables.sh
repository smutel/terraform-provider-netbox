#!/bin/bash

TMP_JSON_FILE="/tmp/docker.json"

if ! which jq > /dev/null; then
  echo "This script need jq binary to work"
fi

curl -sq --unix-socket /var/run/docker.sock http://localhost/containers/json -o ${TMP_JSON_FILE}

DOCKER_COUNT=$(cat ${TMP_JSON_FILE} | jq ".[] | length" | uniq)

for ((i = 0 ; i < ${DOCKER_COUNT} ; i++)); do
  IMAGE_NAME=$(cat /tmp/docker.json | jq .[$i].Image | xargs | cut -d":" -f1)

  if [ "${IMAGE_NAME}" == "netboxcommunity/netbox" ]; then
    PORTS_COUNT=$(cat /tmp/docker.json | jq ".[$i].Ports | length")

    for ((j = 0 ; j < ${PORTS_COUNT} ; j++)); do
      PORT=$(cat /tmp/docker.json | jq .[$i].Ports[$j].PublicPort | xargs)
      if [ "$PORT" != "null" ]; then
        PUBLIC_PORT=$PORT
      fi
    done
  fi
done

export NETBOX_URL="127.0.0.1:$PUBLIC_PORT"
export NETBOX_SCHEME="http"
export NETBOX_TOKEN="0123456789abcdef0123456789abcdef01234567"
