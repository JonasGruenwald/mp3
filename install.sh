#!/usr/bin/env bash

# build binary right into bin folder and make executable
go build -o /usr/local/bin/mp3
chmod +"x" /usr/local/bin/mp3
