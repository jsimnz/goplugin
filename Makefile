# Simple makefile for building and testing GoPlugin

# Run the test suite. Build/Compile the test plugin if necessary
test: test_plugin/test_plugin.so
	@echo ">> RUNNING TEST SUITE..."
	go test

# Build/Compile test plugin as shared lib
test_plugin/test_plugin.so: test_plugin/test_plugin.go test_plugin/tp_gen.go
	@echo ">> BUILDING TEST PLUGIN..."
	go build -buildmode="c-shared" -o test_plugin/test_plugin.so test_plugin/*.go 

# Remove compiled assets
clean:
	rm test_plugin/test_plugin.so test_plugin/test_plugin.h