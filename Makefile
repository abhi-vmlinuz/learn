BINARY    := learn
PREFIX    := /usr/local
DESTDIR   :=
BINDIR    := $(PREFIX)/bin
SHELL_NAME := $(notdir $(SHELL))
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS   := -ldflags "-s -w -X main.version=$(VERSION)"
GOFLAGS   :=

.PHONY: all build install uninstall clean test fmt vet completions help

all: build

build:
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) .

install: build
	install -Dm755 $(BINARY) $(DESTDIR)$(BINDIR)/$(BINARY)
	@case "$(SHELL_NAME)" in \
		bash) \
			./$(BINARY) completion bash > $(BINARY).completion; \
			install -Dm644 $(BINARY).completion $(DESTDIR)$(PREFIX)/share/bash-completion/completions/$(BINARY); \
			rm -f $(BINARY).completion; \
			echo "Installed bash completion to $(PREFIX)/share/bash-completion/completions/$(BINARY)";; \
		zsh) \
			./$(BINARY) completion zsh > $(BINARY).completion; \
			install -Dm644 $(BINARY).completion $(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BINARY); \
			rm -f $(BINARY).completion; \
			echo "Installed zsh completion to $(PREFIX)/share/zsh/site-functions/_$(BINARY)";; \
		fish) \
			./$(BINARY) completion fish > $(BINARY).completion; \
			install -Dm644 $(BINARY).completion $(DESTDIR)$(PREFIX)/share/fish/vendor_completions.d/$(BINARY).fish; \
			rm -f $(BINARY).completion; \
			echo "Installed fish completion to $(PREFIX)/share/fish/vendor_completions.d/$(BINARY).fish";; \
		*) \
			echo "Unknown shell: $(SHELL_NAME). Skipping completion install.";; \
	esac

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY)
	@case "$(SHELL_NAME)" in \
		bash) rm -f $(DESTDIR)$(PREFIX)/share/bash-completion/completions/$(BINARY);; \
		zsh)  rm -f $(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BINARY);; \
		fish) rm -f $(DESTDIR)$(PREFIX)/share/fish/vendor_completions.d/$(BINARY).fish;; \
	esac

clean:
	rm -f $(BINARY)
	rm -f bash-completion zsh-completion fish-completion

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

completions: completions-bash completions-zsh completions-fish

completions-bash: build
	./$(BINARY) completion bash > bash-completion

completions-zsh: build
	./$(BINARY) completion zsh > zsh-completion

completions-fish: build
	./$(BINARY) completion fish > fish-completion

help:
	@echo "Targets:"
	@echo "  build           Build the binary"
	@echo "  install         Install binary to $(DESTDIR)$(BINDIR) and bash completion"
	@echo "  uninstall       Remove installed binary and completion"
	@echo "  clean           Remove build artifacts"
	@echo "  test            Run tests"
	@echo "  fmt             Format Go source"
	@echo "  vet             Run go vet"
	@echo "  completions     Generate all shell completion scripts"
	@echo "  help            Show this help"
