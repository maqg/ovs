all::
	go  build -o ovs main.go
	go  build -o ovsboot ovsboot.go

clean::
	rm -rf ovs ovsboot

fmt::
	go fmt .

package::
	./PACKAGE.sh

build::
	./PACKAGE.sh
