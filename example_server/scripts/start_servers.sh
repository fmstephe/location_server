#!/bin/bash

cd ..
example_server 2> scripts/example.log &
loc_server -m 2> scripts/loc.log &
msg_server 2> scripts/msg.log &
cd -
