![Discord RPG Bot](https://media.giphy.com/media/2h4ObuHvaYnOaHYmnM/giphy.gif)

# Discord RPG Bot

RPG bot for Discord made in go, with some rules I created a few years ago based on Warhammer Fantasy Roleplay V2

## Features

- Character creation with randomized stats
- Proceduraly generated Dungeon exploration with dangerous enemies and traps
- Turn based Dueling between players and foes
- Real Time visualization of the map on browser

## How to play

- Write this in Discord to show the different commands available
  `-crpg Help`

## For Developers

### > Understand the project

For addittionnal informations see [The Documentation](https://drive.google.com/drive/folders/1EmS7LJcMxZhxdygZKR0iP2USDl01H7QP?usp=sharing)

### > Local installation

If you use docker you will only need:

- Docker;
- Docker-Compose;

Refer to [Docker-Setup](#docker-setup) to install with docker.

To run this project, you will also need to install the following dependencies on your system:

- [go](https://golang.org/doc/install)

### > How to launch the project

- To run the project
  `docker-compose up --build`

- On another terminal, after the docker is up
  `docker-compose exec go /bin/sh`
  `go run main.go`

## Contributing

Contributions are always welcome!

See `CONTRIBUTING.md` for ways to get started.

Please adhere to this project's `code of conduct`.
