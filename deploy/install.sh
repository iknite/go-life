#!/bin/sh
set -eu

cmd_name="raft-life"

platform="$(uname -s)"
case "${platform}" in
    Linux*)     machine=linux;;
    Darwin*)    machine=darwin;;
    *)          echo "UNKNOWN platform" && exit 1
esac

arch="$(uname -m)"
case "${arch}" in
    armv6l)    arch=armv6;;
    armv7l)    arch=armv7;;
    *)         arch="amd64"
esac

test -n "${VERSION}" \
    && version=${VERSION} \
    || version="0.0.1"

bin_url="https://github.com/iknite/raft-life/releases/download/${version}/${cmd_name}_${version}_${machine}_${arch}"

bin_dir= $(mktemp -d /tmp/${cmd_name}_XX${bin_dir}XXXXX)
mkdir -m a=rwx -p ${bin_dir}
cd ${bin_dir} && curl -s -L ${bin_url} -O

cmd=${cmd_name}_${version}_${machine}_${arch}
chmod u=rwx ${bin_dir}/${cmd}

test -n "${GOPATH}" \
    && mv ${bin_dir}/${cmd} ${GOPATH}/${cmd_name} \
    || mv ${bin_dir}/${cmd} ${PATH}/${cmd_name}

rm -rf ${bin_dir}
