#!/bin/sh

DISTDIR=package

rm -rf $DISTDIR
mkdir $DISTDIR

TODAY=$(date +"%m%d%Y")

echo "Building..."
go build -o $DISTDIR/ovs main.go
go build -o $DISTDIR/ovsboot ovsboot.go

cd $DISTDIR

mkdir -p logs

if [ ! -x ovs ]; then
	echo "Build ovs error"
	exit 1
fi

if [ ! -x ovsboot ]; then
	echo "Build ovsboot error"
	exit 1
fi

echo "Stripping binaries..."
strip ovs
strip ovsboot

PACKAGE=ovs-$TODAY".tgz"

echo "Packaging..."
cp ../config.yml config-raw.yml
cp ../restart.sh .
cp ../install.sh .
tar -zcf $PACKAGE ../frontend ovs ovsboot logs *.yml *.sh

echo ""
echo "Build Successfully, got file $PACKAGE"
echo ""
