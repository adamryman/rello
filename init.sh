#!/bin/bash

if [ ! -f ./db/db.db ]; then
	sqlite ./db/db.db < ./db/db.sql	
fi
