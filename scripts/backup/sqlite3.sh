#!/usr/bin/env bash

cd ../..
time="$(date +\%Y\%m\%d\%H\%M\%S)"

sqlite3 yitu.db  << EOF
.backup ${time}.db
EOF
