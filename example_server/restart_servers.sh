#!/bin/bash

./kill_servers.sh
git fetch origin
git rebase origin/master
gb
./start_servers.sh
