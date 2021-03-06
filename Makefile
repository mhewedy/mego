all: clean build

build:
	# build the ui
	npm --prefix mego-ui install mego-ui
	cd mego-ui && npm run build && cd -
	mkdir -p mego-api/public
	mv mego-ui/dist/* mego-api/public/
	#build the api
	-cd mego-api; ${GOPATH}/bin/pkger 2>/dev/null; ${GOPATH}/bin/pkger
	cd mego-api;  go build $(cflags)

clean:
	rm -rf mego-ui/dist/
	rm -rf mego-api/public/
	rm -rf mego-api/pkged.go
	rm -rf mego-api/mego
