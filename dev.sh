#!/bin/bash

set -e    # Die if any command fails

MODE="dev"
APP_NAME=emrs
DOCKER_FILE_DEV=docker/Dockerfile
DOCKER_FILE_REL=docker/Dockerfile.rel
DOCKER_TARGET=$DOCKER_FILE_DEV
UNIT_TEST_LOCATION=./scripts/unit-tests.sh

function doUsage() {
  echo -e "\nDeveloper script\n"  
  echo " -h | --help                Help"
  echo " -t | --test                Execute unit tests"
  echo "dev | rel                   Set dev, vs release mode (default:dev)"
  echo "build                       Build the specified mode's docker container"
  echo "run                         Run the specified mode's docker container"
  echo -e "\n"
  echo "example:    ./dev.sh rel build run        Build and run the release mode"
  echo "example:    ./dev.sh dev build            Build the dev mode container"
  echo "example:    ./dev.sh dev run              Run the dev mode container"
  echo "example:    ./dev.sh run                  Run the dev mode container, with dev from default"
  echo "example:    ./dev.sh --test               Just run the unit tests (not containered)"
  echo -e "\n"
}

function doDockerBuild() {
  echo "build:" "$APP_NAME:$MODE" " using " $DOCKER_TARGET
  docker build --tag "$APP_NAME:$MODE" -f "$DOCKER_TARGET" .
}

function doLaunch() {
  echo "run:" "$APP_NAME:$MODE"
  docker run --publish 8080:8080 "$APP_NAME:$MODE"
}

for i in "$@"; do
  case $i in
    run)
      doLaunch
      shift
      ;;
    build)
      doDockerBuild
      shift
      ;;
    rel)
      MODE="rel"
      DOCKER_TARGET=$DOCKER_FILE_REL
      shift
      ;;
    dev)
      MODE="dev"
      DOCKER_TARGET=$DOCKER_FILE_DEV
      shift
      ;;
    -t|--test)
      bash $UNIT_TEST_LOCATION 
      shift
      ;;
    -h|--help)
      doUsage
      exit 0
      ;;
    -*|--*)
      doUsage
      exit 1
      ;;
    *)
      ;;
  esac
done


