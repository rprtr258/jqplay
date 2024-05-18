# jqplay

[![OpenCollective](https://opencollective.com/jqplay/backers/badge.svg)](#backers) [![OpenCollective](https://opencollective.com/jqplay/sponsors/badge.svg)](#sponsors)

[jqplay](https://jqplay.org) is a playground for [jq](https://github.com/jqlang/jq). Please put it into good use.

## Development

To develop `jqplay`, you need to have a [Go development environment](http://golang.org/doc/install).
You also need to have Node & Postgresql installed.

### make start

This script will build and start the `jqplay` server with `docker-compose`.

Your Docker needs to have the [buildx](https://docs.docker.com/engine/reference/commandline/buildx/) command.

Point your browser to [`http://localhost:8080/`](http://localhost:8080/).
