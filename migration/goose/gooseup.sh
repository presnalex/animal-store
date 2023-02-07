#!/bin/bash

DBSTRING="user=postgres password=password dbname=postgres host=localhost port=5438  sslmode=disable"

goose postgres "$DBSTRING" up