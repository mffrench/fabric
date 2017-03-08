#!/bin/sh
docker ps -a | grep peer | awk '{print $1}' | xargs docker stop
docker ps -a | grep peer | awk '{print $1}' | xargs docker rm
docker ps -a | grep orderer | awk '{print $1}' | xargs docker rm
docker network ls | grep _default | awk '{print $1}' | xargs docker network rm
cd $GOPATH/src/github.com/hyperledger/fabric/
make clean
make peer; make orderer; make configtxgen; make peer-docker; make orderer-docker
cd -
