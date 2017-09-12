#!/bin/sh

DISTDIR=package

rm -rf $DISTDIR
mkdir $DISTDIR

TODAY=$(date +"%m%d%Y")

echo "Building..."
go build -o $DISTDIR/ovs main.go

cd $DISTDIR

mkdir -p logs

if [ ! -x ovs ]; then
	echo "Build ovs error"
	exit 1
fi

echo "Stripping binaries..."
strip ovs

PACKAGE=ovs-$TODAY".tgz"

echo "Packaging..."
cp ../config.yml config-raw.yml
cp ../restart.sh .
tar -zcf $PACKAGE ../frontend ovs logs *.yml *.sh

echo ""
echo "Build Successfully, got file $PACKAGE"
echo ""
