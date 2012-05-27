#!/bin/bash

cd ..
example -port $1 2> scripts/example.log &
cd -
