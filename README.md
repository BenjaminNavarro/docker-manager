# Docker Manager

## Introduction

With docker manager you can:
 * Configure your containers easily using a simple YAML file
 * Start them in the background (detached mode)
 * Make multiple connections to your running containers
 * Save the containers' state at any time
 * Stop your containers and, optionally, save their current state

## Installation

To install docker manager in your Go workspace, run:

```
go get github.com/BenjaminNavarro/docker-manager
```

## Usage

### The configuration file

Docker manager requires a YAML configuration file located at:
```
$HOME/.docker-manager.yaml
```
It can contain any number of docker configurations. You can follow the following example to write your own configuration file:
```YAML
-   name: MyFirstContainer
    image: my/container
    tag: latest # default
    save_tag: latest    # default = tag
    autosave: true      # default: false
    runtime: nvidia     # default: none
    privileged: true    # default: false
    gui: true           # default: false
    shell: bash         # default
    folders:            # default: none
        -   host: /path/to/my/folder
            container: /path/inside/the/container
    capabilities:
        add:
            - NET_RAW
        drop:
            - SETPCAP
    extra_flags: --env="GOPATH=/opt/go"

-   name: Another Container    # spaces will be removed
    image: other/container
    autosave: true
    network: host               # default bridge
    capabilities:
        add: [ALL]

```

### Using docker manager

You can run `docker-manager` then go through the menus to select what you want to do:

```
Please select an image:
	1) MyFirstContainer
	2) AnotherContainer
>>> 1  
Please select an action for this container:
	1) Start
	2) Connect
	3) Save
	4) Stop
>>> 1
```

You can also pass arguments to `docker-manager` to skip the menus. The first argument is the container's name (as given in `.docker-manager.yaml`) and the second is the action to take:
 * start
 * connect
 * save
 * stop

So for example you can run: `docker-manager MyFirstContainer start`

If one argument is missing or incorrect, the corresponding menu will be prompted.
