# saido
saido means monitor in [Hausa](https://en.wikipedia.org/wiki/Hausa_language)

![Logo](assets/Saido300.jpg)

[![Test-Linux](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml)
[![Test-MacOs](https://github.com/bisohns/saido/actions/workflows/test-macos.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-macos.yml)
[![Test-Windows](https://github.com/bisohns/saido/actions/workflows/test-windows.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-windows.yml)

## Installation
NOTE: We've currently only tested on `Mac Os`, `Windows 11` and `Linux Ubuntu 20.04 LTS`
### Installing Binary
Download latest binary for the target machine from [Github Releases](https://github.com/bisohns/saido/releases/latest)
#### Windows Cmd
```cmd
REM Check your `PATH` 
path
REM Copy the saido binary to a directory that is on your `PATH`
copy <path_to_saido_binary> <path>
```
#### Windows Powershell
```powershell
#Appends to existing path
$env:Path += ";<path_to_saido_binary>" 
```
#### Linux or Mac OS 
```bash 
#Add saido to your `.bashrc` or `.zshrc`
alias saido='<path_to_saido_binary>' 

#Or Copy binary to your bin directory
cp <path_to_saido_binary> /usr/local/bin 
```
### Development
NOTE: `!windows` flag is our current specification of what `unix` means, see [issue](https://github.com/golang/go/issues/20322) for why *_unix.go files will still attempt to run on windows.
#### Requirements
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Yarn](https://classic.yarnpkg.com/lang/en/docs/install/)
- [Air](https://github.com/cosmtrek/air)
```bash
git clone https://github.com/bisohns/saido
cd saido
## Update dependencies
make dependencies

# Build and serve frontend
make app

# Modify generated `config-test.yaml` file and air would reload server
```

## Usage
The simplest usage is running
```bash
# binary is downloaded and named as saido
saido --config <path_to_configuration_yaml_file> --port 3000 --verbose
```
Saido cli flags & commands
```bash
Saido
Tool for monitoring metrics

Usage:
  saido [flags]
  saido [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Get saido version

Flags:
      --config string   Path to config file
  -h, --help            help for saido
  -b, --open-browser    Prompt open browser
  -p, --port string     Port to run application server on (default "3000")
  -v, --verbose         Run saido in verbose mode

Use "saido [command] --help" for more information about a command.
```

## Yaml Configuration File
NOTE: Use single qoutes (`''`) for any string within the config file.
### Hosts
`hosts`
#### Supported hosts command
NOTE: The broad assumptions we make is that the host names are treated as unique.
* `children` - a list of hosts with host names
#### Setting up local connection to host
```yaml
hosts:
  children:
    # Host with host name `localhost`
    'localhost':
        connection:
            # define `connection` with  `type` as `local` 
            type: local
        metrics:
            memory:
            disk:
            tcp:
            docker:
            uptime:
poll-interval: 10
```
#### Setting up ssh connection to host (with ssh username and password)
```yaml
hosts:
   children:
   # Host with host name `192.168.166.167`
    '192.168.166.167':
      connection:
        # define `connection` with `type` as `ssh` 
        type: ssh
        # define `username` and `password` for connection to ssh server
        username: <username>
        password: <password>
      metrics:
        custom-ls: 'ls $HOME/app'
metrics:
  memory:
  disk:
  tcp:
  docker:
  uptime:
poll-interval: 10
```
#### Setting up ssh connection to host (with ssh private key and port)
```yaml
hosts:
   children:
   # Host with host name `0.0.0.0`
     "0.0.0.0":
      connection:
        # define `connection` with `type` as `ssh` 
        type: ssh
        # define `username` and `port` and `private_key_path` for connection to ssh server
        username: <username>
        port: 2222
        private_key_path: "path_to_private_key"
metrics:
  memory:
  disk:
  tcp:
  docker:
  uptime:
poll-interval: 10
```
#### Setting up ssh connection to host (with ssh a password protected private key)
```yaml
hosts:
   children:
   # Host with host name `0.0.0.0`
     "0.0.0.0":
      connection:
        # define `connection` with `type` as `ssh` 
        type: ssh
        # define `username` and `port` , `private_key_path` and `private_key_passphrase` for connection to ssh server
        username: <username>
        port: 2222
        private_key_path: <"path_to_private_key">
        private_key_passphrase: <private_key_passphrase>
metrics:
  memory:
  disk:
  tcp:
  docker:
  uptime:
poll-interval: 10
```
### Metrics
`metrics`
#### Supported metrics command
NOTE: Custom metrics with commands can be added to the configuration file
* `memory` - for calculating memory usage
* `disk`- for calculating disk usage
* `tcp` - for getting tcp connection information
* `docker` - for getting docker container information
* `uptime` - for calculating uptime and idle time of the host
#### Setting Global metrics 
```yaml
hosts:
  children:
    'localhost':
        connection:
            type: local
    '0.0.0.0':
      connection:
        type: ssh
        username: <username>
        port: 2222
        private_key_path: '<path_to_private_key>'
# Metrics are defined on a global scope. Same level with the hosts
metrics:
    memory:
    disk:
    tcp:
    docker:
    uptime:
poll-interval: 10
```
#### Setting Local metrics
```yaml
hosts:
  children:
    'localhost':
        connection:
            type: local
        # Metrics are defined within a host (localhost) scope
        metrics:
            memory:
            disk:
            tcp:
            docker:
            uptime:
poll-interval: 10
```
#### Setting Custom metrics
```yaml
hosts:
  children:
    'localhost':
        connection:
            type: local
        # A custom metric named `custom-ls` with a command `ls $HOME/app` is defined within a host (localhost) scope
        metrics:
            custom-ls: 'ls $HOME/app'   
poll-interval: 10
```
### Polling
`polling-interval` - interval in seconds between requests to host (value must be greater than or equal to 5 seconds)
#### Example
NOTE: Use a reasonable time interval between 10-30 seconds to avoid overloading server
```yaml
hosts:
  children:
    'localhost':
        connection:
            type: local
        metrics:
            disk: 
# Polling interval set to 5 seconds
poll-interval: 5  
```

## Deployment
### Tagging
To create a new tag, use the make file
```bash
make upgrade version=0.x.x
```
This will tag a new version, push to github and trigger goreleaser

## License
The project is opened under the [Apache License](https://github.com/bisohns/saido/blob/master/LICENSE)
### Credits
 - Logo by [Williams Praise](https://github.com/kubyruby)
 - Goph by [Melbahja](https://github.com/melbahja/goph)