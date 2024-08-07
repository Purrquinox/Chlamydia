BUILDTIME := $(shell date '+%Y-%m-%d | %H:%M:%S')
REPONAME := github.com/Purrquinox/Chlamydia
PROJECTNAME := chlamydia

COMBOS := linux/amd64 linux/arm64 windows/amd64 windows/arm64

# Cross-compilers for each platform
CC_linux_amd64 := x86_64-linux-gnu-gcc
CC_linux_arm64 := aarch64-linux-gnu-gcc
CC_windows_amd64 := x86_64-w64-mingw32-gcc
CC_windows_arm64 := aarch64-w64-mingw32-gcc

.PHONY: all release fmt clean prerelease

all:
	@echo "Building in debug mode..."
	CGO_ENABLED=1 go build -v $(GOFLAGS_DBG)

release:
	@python build_release.py

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin

prerelease:
	@echo "Installing prerequisites..."
	sudo apt-get update
	sudo apt-get install -y gcc-multilib g++-multilib \
		x86_64-linux-gnu-gcc \
		aarch64-linux-gnu-gcc \
		x86_64-w64-mingw32-gcc \
		aarch64-w64-mingw32-gcc
