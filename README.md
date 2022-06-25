# IPby

### Simply private/public IP check

## 1. Support OS
- macOS 
- Linux

## 2. How to use it?
This is CLI tool.
1. Quickly and simply use
```bash
ipby
```
2. Advanced options 
```bash
ipby <command>
```
- version: check ipby version
- all: show detail IP for public and private
- public: show detail IP for IP public
- private: show detail IP for private
- help: command explain

## 3. How to own build
Before this, need GIT and GoLang
```bash
git clone https://github.com/leelsey/ipby
cd IPby/cmd/ipby &&  go mod tidy
go build ipby.go
```

## 4. OpenSource
This is MIT License.
Public IP searched by [ipify](https://www.ipify.org/).
Private IP searched on ipconfig for macOS and hostname for Linux.
