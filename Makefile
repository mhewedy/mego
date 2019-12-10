
build:
	# build the ui
	cd mego-ui && npm run build
	mkdir mego-api/public
	mv mego-ui/dist/* mego-api/public/

clean:
	rm -rf mego-ui/dist/
	rm -rf mego-api/public/
