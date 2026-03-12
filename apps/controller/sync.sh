#!/bin/bash

source_dir=/home/rawad/eigen-db/apps/controller # on personal Arch Linux machine
target_dir=/home/rawad/eigen-db/apps # on FreeBSD development server

ssh rawad@freebsd "rm -rf /home/rawad/eigen-db/apps/controller"
scp -r $source_dir rawad@freebsd:$target_dir