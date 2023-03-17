# PMU

`pmu` is a small tool that offers a CLI similar to that of [pm2](https://github.com/Unitech/pm2),
but instead of running a daemon to manage processes, it just creates systemd services files,
and forwards commands to systemctl and journalctl.

It provides the ease of use offered by pm2 and the ubiquity and reliability of systemd
without the need to run any extra node-specific software in the background

## Why?

For when you want to use systemd systems where it's already installed, but also want the convenience of being able to
spin up a new service in a single command without having to mess with configurations.


## Installation
`pmu` is written in Go and distributed as a single binary executable.
1. Move the executable to the desired directory
2. Add this directory to the PATH environment variable
3. Verify that you have execute permission on the file


## Usage

## Notes

Simplest way to daemonize a node app
```shell
pmu start app.js
```
Starting other application types
```shell
pmu start bashscript.sh
pmu start python-app.py
pmu start ./binary-file -- --port 1520
```
For python and node scripts a default interpreter will be set,
but you can also specify it with the `--interpreter` flag.

* I don't like the default behaviour of pm2 to not start apps on startup, so by default pmu apps are "enabled" with systemctl
* Unlike the default behaviour of systemctl services though, pmu apps are set to always restart automatically on exit 
