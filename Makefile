# Setup name variables for the package/tool
NAME := fileserver
PKG := github.com/urvil38/$(NAME)

CGO_ENABLED := 0

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild: