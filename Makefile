.PHONY: bump-dependencies
bump-dependencies:
	go get common@main
	go mod tidy
	cd testsuite;\
		go get common@main storx@main uplink@main;\
		go mod tidy;\