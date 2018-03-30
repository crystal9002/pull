#!/bin/bash
cd $1
if [ $2 != "" ]; then
  git checkout master
  if [ $2 != "master" ]; then
    git branch -D $2
    git checkout $2
  fi
fi
git pull
