# Discord RPG Bot

RPG bot for Discord made in go, with some rules I created a few years ago based on Warhammer Fantasy Roleplay V2

## Requirements

If you use docker you will only need:

- Docker;
- Docker-Compose;

Refer to [Docker-Setup](#docker-setup) to install with docker.

To run this project, you will also need to install the following dependencies on your system:

- [go](https://golang.org/doc/install)

## How to launch the project

- To run the project
  `docker-compose up --build`

- On another terminal, after the docker is up
  `docker-compose exec go /bin/sh`
  `go run main.go`

## Help command

- Write this in Discord to show the different commands available
  `-crrpg Help`
