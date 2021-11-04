[![CircleCI](https://circleci.com/gh/Sharpz7/sharpcd.svg?style=svg)](https://circleci.com/gh/Sharpz7/sharpcd)

# SharpCD || Continuous Development for your server

![](https://files.mcaq.me/zbnf.png)

# Example Config
```yml
version: 1

trak:
  local: https://localhost:5666
  remote: https://mcaq.me:5666

tasks:
  basic_task:
    name: Basic
    type: docker
    sharpurl: https://localhost:5666
    giturl: https://raw.githubusercontent.com/Sharpz7/
    compose: /sharpcd/dev/testing/basic.yml

  registry_task:
    name: Registry
    type: docker
    registry: docker.mcaq.me
    envfile: .env
    sharpurl: https://localhost:5666
    giturl: https://raw.githubusercontent.com/Sharpz7/
    compose: /sharpcd/dev/testing/registry.yml

  env_task:
    name: Enviroment Test Fail
    type: docker
    envfile: .env
    sharpurl: https://localhost:5666
    giturl: https://raw.githubusercontent.com/Sharpz7/
    compose: /sharpcd/dev/testing/env.yml
```

# Installation
On linux, just run:
```console
╭─adam@box ~/
╰─➤  sudo curl -s -L https://github.com/Sharpz7/sharpcd/releases/download/3.1/install.sh | sudo bash
```

Or for just the client:
```console
╭─adam@box ~/
╰─➤  sudo curl -s -L https://github.com/Sharpz7/sharpcd/releases/download/3.1/install.sh | sudo bash -s client
```

## Command Options

On linux, just run:
```console
╭─adam@box ~/
╰─➤  sharpcd help

Args of SharpCD:

        - server: Run the sharpcd server
        - setsecret: Set the secret for API and Task Calls
        - addfilter: Add a url for a compose file
        - changetoken: Add a token for private github repos
        - removefilter: Remove a url for a compose file
        - version: Returns the Current Version
        - trak: Run the Trak program

Sub Command Trak:

        - alljobs {type}: Get info on all jobs
        - job {type} {id}: Get info on job with logging
        - list {type}: Get all jobs running on sharpcd server

Flags:

  -secret string
        Put secret as a arg for automation tasks
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)