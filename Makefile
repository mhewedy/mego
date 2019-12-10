all: clean build

build:
	# build the ui
	cd mego-ui && npm run build
	mkdir -p mego-api/public
	mv mego-ui/dist/* mego-api/public/
	cd mego-api && pkger && go build

clean:
	rm -rf mego-ui/dist/
	rm -rf mego-api/public/
	rm -rf mego-api/pkged.go
	rm -rf mego-api/mego
