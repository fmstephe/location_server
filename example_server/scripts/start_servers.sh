#!/bin/bash

cd ..
example_server > scripts/example.log &
loc_server -m > scripts/loc.log &
msg_server > scripts/msg.log &
cd -
