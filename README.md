# saido
saido means monitor in Hausa


![Logo](assets/Saido300.jpg)

[![Test-Linux](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-ssh.yml)
[![Test-MacOs](https://github.com/bisohns/saido/actions/workflows/test-macos.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-macos.yml)
[![Test-Windows](https://github.com/bisohns/saido/actions/workflows/test-windows.yml/badge.svg)](https://github.com/bisohns/saido/actions/workflows/test-windows.yml)

NOTE: `!windows` flag is our current specification of what `unix` means, seee [issue](https://github.com/golang/go/issues/20322) for why *_unix.go files will still attempt to run on windows


## Installation

### Usage

For personal usage, install latest from [Github Releases](https://github.com/bisohns/saido/releases) 

```bash
# binary is downloaded and named as saido
saido api
```


### Development

With Golang installed, run

```bash
git clone https://github.com/bisohns/saido
cd saido
## Update Golang dependencies
go get .

## Update yarn dependencies
cd web
yarn install

# Run websocket server and serve frontend
go run main.go api
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
