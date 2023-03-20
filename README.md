# MP3

![Header image (an mp3 player)](header.png)

MP3 is a small tool that offers a CLI similar to that of [pm2](https://github.com/Unitech/pm2),
but instead of running a daemon to manage processes, it just creates systemd services files,
and forwards commands to systemctl and journalctl.

It provides the ease of use offered by pm2 and the ubiquity and reliability of systemd
without the need to run any extra node-specific software in the background

The name is purely to create confusion.

## Motivation

I want to use systemd on ubuntu where it's already installed, but also want the convenience of being able to
spin up a new service in a single command without having to mess with configurations.

I also want a filtered overview of systemd services that I created myself, including a nice way to display their status
and logs.

## Installation

mp3 is written in Go and can be distributed as a single binary executable.

1. Move the executable to the desired directory
2. Add this directory to the PATH environment variable
3. Verify that you have execute permission on the file

## Usage

### Start an app

Simplest way to daemonize a node app (just like pm2)

```shell
mp3 start app.js
```

Starting other application types

```shell
mp3 start bashscript.sh
mp3 start python-app.py
mp3 start ./binary-file -- --port 1520
```

For python and node scripts a default interpreter will be set,
but you can also specify it with the `--interpreter` flag.

pm2 compatible flags for `start`:

```shell
# Specify a name for your new app
--name <app_name>

# Pass extra args
-- arg1 arg2 arg3

# Delay between automatic restarts
--restart-delay <delay in ms>

# Do not auto restart app
--no-autorestart
```

Unlike pm2, mp3 will enable apps when started, so they get automatically started on startup as as well. To disable this
behaviour, you can use the `--create-only` flag described below

mp3 specific flags

```shell
# Only create a service file, do not start and enable
--create-only
```

After an application has been started once, you can always start it again from anywhere with `mp3 start <name>`.


### Managing processes

Commands are the same as in pm2, but just call systemctl commands with the full service name

```shell
mp3 restart app_name
mp3 reload app_name
mp3 stop app_name
mp3 delete app_name
```

For all commands except delete you can also pass 'all' instead of an app name.
### Status & Logs

You can display the status of all mp3 services

```shell
mp3 [list|ls|status]
```

Display the logs for all mp3 services

```shell
mp3 logs
```

or for a specific app

```shell
mp3 logs <app_name>
```

By default, logs will show the last 100 lines and tail. If you want to change the way logs are displayed, you can pass
on any arguments to `journalct` with `--`.

If you pass any arguments using this method, the defaults (100 lines tailing) will not be used.
