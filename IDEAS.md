# Connect
Special Code for starting caddy service with its recommended service file, and setting up a caddyfile structure that 
allows dynamically adding sites by creating new files in a folder, files are automatically imported caddyfiles.
```shell
mp3 start caddy
```

magic command to connect a mp3 service to a domain, by detecting its port and creating a caddyfile config for it

```go
mp3 connect app app.example.com
```

# Notifications

Add option to create a service which all managed services call with  OnFailure=

The services just executes mp3 with a flag including the failed service's name and mp3 sends out a notification via Telegram.

Required parameters like a tegram token and chat id can be stored directly in the notification service file's environment variables
