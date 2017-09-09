all::
	go build -o ovs main.go

clean::
	rm -rf ovs

fmt::
	go fmt .

package::
	./PACKAGE.sh

build::
	./PACKAGE.sh
