# Adopt

`mp3 adopt service_path` Adopt a regular systemd service into the mp3 namespace, making it show up in status and logs.

This will just [create a symlink to the units](https://serverfault.com/a/1078481/1010774) service file, e.g. for 
`mongod.service` it will create a symlink called `mp3.mongod.service` which will be an alias for `mongod` and cause it 
to show up in mp3.

`mp3 delete mongod` will of course only delete  the symlink and leave the original service untouched.

# Connect

Special Code for starting caddy service with its recommended service file, and setting up a caddyfile structure that
allows dynamically adding sites by creating new files in a folder, files are automatically imported caddyfiles.

```shell
mp3 setup caddy
```

magic command to connect a mp3 service to a domain, by detecting its port and creating a caddyfile config for it

```go
mp3 connect app app.example.com
```

Use: https://pkg.go.dev/github.com/cakturk/go-netstat/netstat?utm_source=godoc

# Notifications

Add option to create a service which all managed services call with OnFailure=

The services just executes mp3 with a flag including the failed service's name and mp3 sends out a notification via
Telegram.

Required parameters like a tegram token and chat id can be stored directly in the notification service file's
environment variables
