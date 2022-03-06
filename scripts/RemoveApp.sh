#!/bin/bash

kill -s SIGTERM $(pgrep -f go);

rm -rf /var/app/current;