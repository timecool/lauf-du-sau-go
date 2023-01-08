#!/bin/bash
cowsay LAUF DU SAU
go build -o main
systemctl stop goapp
systemctl start goapp