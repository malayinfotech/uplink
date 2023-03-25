#!/usr/bin/env bash
set -ueo pipefail
set +x

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $SCRIPTDIR

# setup tmpdir for testfiles and cleanup
TMP=$(mktemp -d -t tmp.XXXXXXXXXX)
cleanup(){
	rm -rf "$TMP"
}
trap cleanup EXIT

STORX_GO_MOD=${STORX_GO_MOD:-"../go.mod"}
VERSION=$(go list -modfile $STORX_GO_MOD -m -f "{{.Version}}" storx)


go install storx/cmd/certificates@$VERSION
go install storx/cmd/identity@$VERSION
go install storx/cmd/satellite@$VERSION
go install storx/cmd/storagenode@$VERSION
go install storx/cmd/versioncontrol@$VERSION
go install storx/cmd/storx-sim@$VERSION
go install storx/cmd/multinode@$VERSION
go install gateway@latest

echo "Used version: $VERSION"

export STORX_NETWORK_DIR=$TMP

STORX_NETWORK_HOST4=${STORX_NETWORK_HOST4:-127.0.0.1}
STORX_SIM_POSTGRES=${STORX_SIM_POSTGRES:-""}

# setup the network
# if postgres connection string is set as STORX_SIM_POSTGRES then use that for testing
if [ -z ${STORX_SIM_POSTGRES} ]; then
	storx-sim -x --satellites 1 --host $STORX_NETWORK_HOST4 network setup
else
	storx-sim -x --satellites 1 --host $STORX_NETWORK_HOST4 network --postgres=$STORX_SIM_POSTGRES setup
fi

sed -i 's/# metainfo.rate-limiter.enabled: true/metainfo.rate-limiter.enabled: false/g' $(storx-sim network env SATELLITE_0_DIR)/config.yaml

storx-sim -x --satellites 1 --host $STORX_NETWORK_HOST4 network test bash "$SCRIPTDIR"/rclone.sh
storx-sim -x --satellites 1 --host $STORX_NETWORK_HOST4 network destroy
