#!/bin/bash

./kill_servers.sh
git fetch origin
git rebase origin/master
cd ../
gb
cd -
./start_servers.sh
