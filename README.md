[![CircleCI](https://circleci.com/gh/SharpSet/sharpcd.svg?style=svg)](https://circleci.com/gh/SharpSet/sharpcd)

![SharpCD](https://files.mcaq.me/043h0.png)
# Continuous Development for your server!

SharpCD is a simple, yet powerful, continuous development tool for your server.
It allows you to easily deploy, manage and track docker-compose projects from any location.

![](https://files.mcaq.me/r4844.png)

# Example Config

Everything is controlled by the `sharpcd.yml` file. You can also use
remote configuration files, allowing for easy deployment of multiple projects that can
have multiple dependencies.

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
    giturl: https://raw.githubusercontent.com/SharpSet/
    compose: /sharpcd/dev/testing/basic.yml

  registry_task:
    name: Registry
    type: docker
    registry: docker.mcaq.me
    envfile: .env
    sharpurl: https://localhost:5666
    giturl: https://raw.githubusercontent.com/SharpSet/
    compose: /sharpcd/dev/testing/registry.yml

  env_task:
    name: Enviroment Test
    type: docker
    envfile: .env
    sharpurl: https://localhost:5666
    giturl: https://raw.githubusercontent.com/SharpSet/
    compose: /sharpcd/dev/testing/env.yml
```

# Installation
On linux, just run:
```console
sudo curl -s -L https://github.com/Sharpz7/sharpcd/releases/download/3.7/install.sh | sudo bash
```

Or for just the client:
```console
sudo curl -s -L https://github.com/Sharpz7/sharpcd/releases/download/3.7/install.sh | sudo bash -s client
```

## Command Options

On linux, just run:
```console
foo@bar:~$ sharpcd help

Args of SharpCD:

        - server: Run the sharpcd server
        - setsecret: Set the secret for API and Task Calls
        - addfilter: Add a url for a compose file
        - changetoken: Add a token for private github repos
        - removefilter: Remove a url for a compose file
        - version: Returns the Current Version
        - trak: Run the Trak program

Sub Command Trak:

        - sharpcd trak alljobs {location}
                Get info on all jobs

        - sharpcd trak job {location} {job_id}
                Get info on job with logging

        - sharpcd trak list {location}
                Get all jobs running on sharpcd server

        - sharpcd trak logs {location} {job_id}
                Get Logs from a Job

Flags:

  -remotefile string
        Location of Remote sharpcd.yml file
  -secret string
        Put secret as a arg for automation tasks
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)