#!/bin/bash

./kill_servers.sh
git fetch origin
git rebase
gb
./start_servers.sh
