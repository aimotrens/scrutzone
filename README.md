# ScrutZone

ScrutZone, a portmanteau word of "Scrutinize" and "Zone" (zone in the sense of “your IT realm that you want to monitor”), is a tool to watch the availability of a service. It is designed to be easy to use compared to other tools like Nagios, Zabbix, Icinga etc.


## Features
- **Simple**: You can start monitoring a service with just a few lines of YAML.
- **Flexible**: You can monitor a service in various ways.
- **Scalable**: You can monitor multiple services with a single ScrutZone instance.


## Installation
Just download the binary from the [release page](https://github.com/aimotrens/scrutzone/releases). ScrutZone is a single binary, so you can run it without any dependencies.

Alternatively, you can use the docker image. The docker image is available at [Docker Hub](https://hub.docker.com/r/t3a6/scrutzone).


## Usage
Create a configuration file in YAML format. Examples are available in the `config.examples` directory.

You must create a yaml file named `scrutzone.yml` in the `./config` directory or define a different config file path with ENV variable `SCRUTZONE_CONFIG_FILE`.

In this file, you must define a config folder (`checkConfigDir`) for the checks. ScrutZone will read all the files in this folder recursively.
