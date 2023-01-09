# Go executable
GO = go

# Build flags
ADDR = 4000

# Executable name
EXECUTABLE = app

# Build the app
build:
	$(GO) build $(BUILD_FLAGS) -o $(EXECUTABLE)

# Cleanup
clean:
	rm -f $(EXECUTABLE)

# Run the app
run:
	$(GO) run ./cmd/web/ -addr $(ADDR)
