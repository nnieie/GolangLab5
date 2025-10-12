#!/bin/bash

(go run ./cmd/api) &
sleep 2s
printf "\n"

(go run ./cmd/user) &
sleep 2s
printf "\n"

wait