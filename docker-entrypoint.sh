#!/usr/bin/env bash

set -o errexit

echo "Docker Entrtypoint"

echo "$@"

case "$1" in
  w)
    set -- ./main -w th_text.txt
  ;;
  m)
    set -- ./main -m th_text.txt
  ;;
  c)
    set -- ./main -c th_text.txt
  ;;
  l)
    set -- ./main -l th_text.txt
  ;;
  *)
    set -- ./main th_text.txt
  ;;
esac

exec "$@"

