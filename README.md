# saido
saido means monitor in [Hausa](https://en.wikipedia.org/wiki/Hausa_language)


![Logo](assets/Saido300.jpg)

[![Test-Linux](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml)
[![Test-MacOs](https://github.com/bisohns/saido/actions/workflows/test-macos.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-macos.yml)
[![Test-Windows](https://github.com/bisohns/saido/actions/workflows/test-windows.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-windows.yml)


## Installation
NOTE: We've currently only tested on `Mac Os`, `Windows 11` and `Linux Ubuntu 20.04 LTS`
### Requirements
- [Golang](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Yarn](https://classic.yarnpkg.com/lang/en/docs/install/)
- [Air](https://github.com/cosmtrek/air)
### Installing Binary
Download latest binary for the target machine from [Github Releases](https://github.com/bisohns/saido/releases/latest)

Windows Cmd
```cmd
REM Check your `PATH` 
path
REM Copy the saido binary to a directory that is on your `PATH`
copy <path_to_saido_binary> <path>
```
Windows Powershell
```powershell
#Appends to existing path
$env:Path += ";<path_to_saido_binary>" 
```
Linux or Mac OS 
```bash 
#Add saido to your `.bashrc` or `.zshrc`
alias saido='<path_to_saido_binary>' 

#Or Copy binary to your bin directory
cp <path_to_saido_binary> /usr/local/bin 
```

### Development

NOTE: `!windows` flag is our current specification of what `unix` means, seee [issue](https://github.com/golang/go/issues/20322) for why *_unix.go files will still attempt to run on windows.

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
saido --config config.yaml --port 3000 --verbose
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

### Yaml Configuration File
NOTE: Use single qoutes (`''`) for any string within the config file. Also, the broad assumptions we make is that the host names are treated as unique.
* Server (localhost)
```yaml

```
* Metrics (Global / Local / Custom)
```yaml

```
* Ssh key files (With and Without password)
```yaml

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
