all: clean build

build:
	# build the ui
	npm --prefix mego-ui install mego-ui
	cd mego-ui && npm run build && cd -
	mkdir -p mego-api/public
	mv mego-ui/dist/* mego-api/public/
	#build the api
	cd mego-api && go build
	ls -lt mego-api

clean:
	rm -rf mego-ui/dist/
	rm -rf mego-api/public/
	rm -rf mego-api/pkged.go
	rm -rf mego-api/mego
