#!/bin/bash

dirpath=$(dirname $(which $0))
cd "$dirpath"/..

echo -n "Choose a profile (dev|test|prod) [dev]: "
read -r profil
if [ -z $profil ]; then
    profil="dev"
fi

echo $profil
