#!/usr/bin/env bash

run_tests() {

    docker-compose -f docker-compose.test.yml up --abort-on-container-exit

}

run_app() {

    docker-compose -f docker-compose.yml up 
}

print_help() {
    echo "URL shortener app helper script"
    echo
    echo "Syntax: run [-r|t|h]"
    echo "-r|--run-app      Run URL shortener application"
    echo "-t|--test         Run test suite"
    echo "-h|--help         Print this help menu"
}

while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
    -t|--test)
        run_tests
        shift
        ;;
    -r|--run-app)
        run_app
        shift
        ;;
    -h|--help|*)
        shift
        print_help
        ;;
  esac
done  
