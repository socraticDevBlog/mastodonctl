#!/bin/source
# shellcheck shell=sh

###################################################################
#Script Name	: setup local development environment
#Description	:
#Args           : None
#Author       	: socraticDev
#Email         	: socraticdev@gmail.com

# IMPORTANT: run with 'source' command  (ex.: source ./setup.sh)
#            in order to set environment variables
###################################################################

echo "Install app dependency on local machine"
go install

echo "Build app into an executable binary"
go build .

if [ "${GOPATH}" ]; then
    echo "GOPATH is set: that's good!"
else
    echo ""
    echo "(warning) GOPATH variable is unset! Please learn how to set your GOPATH and do it!"
    echo ""
fi

echo "Export configuration file (conf.json) file to environment variable"
curr_dir=$(pwd)
export CONFIG_FILEPATH="${curr_dir}"/conf.json
echo "Your config file is located at:"
echo "${CONFIG_FILEPATH}"

echo ""

echo "Now you can run mastodonctl simply by issuing this command:"
echo ""
echo "mastodonctl"
