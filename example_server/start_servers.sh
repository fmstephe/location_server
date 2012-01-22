#!/bin/bash

example_server > example.log &
loc_server -m > loc.log &
msg_server > msg.log &
