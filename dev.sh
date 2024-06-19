#!/bin/bash

set -e    # Die if any command fails

MODE="dev"

APP_NAME=emrs

DOCKER_FILE_DEV=docker/Dockerfile
DOCKER_FILE_REL=docker/Dockerfile.rel
DOCKER_TARGET=$DOCKER_FILE_DEV

TEST_MODULE_LIST=(
  badger
  reaper
)

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

function confirmChoice() {
  read -p "Press [y|Y] to confirm: " -n 1 -r
  echo    # (optional) move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
      echo -e "\n\tconfirmed\n"
      return
  fi
  echo "Exiting. User decided not to continue."
  exit 2
}

function doRemoveAllImages() {
  docker rmi -f $(docker images -aq)
}

function doRemoveAllImagesAndVolumes() {
  docker rm -vf $(docker ps -aq)
}

function doTests() {
  echo "[ RUNNING TESTS ]"

  go clean -cache
  
  for module in ${TEST_MODULE_LIST[*]}; do
    cd ${module}
    go clean -cache
    go test . -v
    cd -
  done
}

for i in "$@"; do
  case $i in
    run)
      doLaunch
      exit 0
      ;;
    build)
      doDockerBuild
      exit 0
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
    clean)
      echo -e "\n\tWARNING:\n\n\tThis will remove ALL docker images\n\n"
      confirmChoice
      doRemoveAllImages
      exit 0
      ;;
    purge)
      echo -e "\n\tWARNING:\n\n\tThis will remove ALL docker IMAGES & VOLUMES\n\n"
      confirmChoice
      doRemoveAllImagesAndVolumes
      exit 0
      ;;
    -t|--test)
      doTests
      exit 0
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


