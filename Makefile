include ./standard.mk

HOMEBREW_PACKAGES=bash coreutils editorconfig-checker findutils git git-cliff git-lfs go grep jq k1LoW/tap/tbls nodejs pre-commit python@3.11 trivy trufflesecurity/trufflehog/trufflehog

#-------------------------------------------------------------------------------
# Environment

.PHONY: install-hooks
## install-hooks: [tools] Install/upgrade the Git hooks used for ensuring consistency.
install-hooks:
	@ $(HEADER) "=====> Installing Git hooks..."
	cp -vf .githooks/commit-msg.sh .git/hooks/commit-msg
	chmod +x .git/hooks/*
	pre-commit install

	@ $(BORDER) "Learn more about 'pre-commit' at:" "  https://pre-commit.com" " " "Learn more about 'gommit' at:" "  https://github.com/antham/gommit"

# goplicate-start:linting
#-------------------------------------------------------------------------------
# Linting

.PHONY: pre-commit
## pre-commit: [lint]* Runs `pre-commit` against all files.
pre-commit:
	@ $(HEADER) "=====> Running pre-commit..."
	pre-commit run --all-files

# goplicate-end:linting

# goplicate-start:git
#-------------------------------------------------------------------------------
# Git Tasks

.PHONY: changelog
## changelog: [release]* Generates the CHANGELOG for the release.
changelog:
	@ $(HEADER) "=====> Updating the CHANGELOG..."
	git cliff -o CHANGELOG.md

.PHONY: tag
## tag: [release]* Signs and tags the release.
tag:
	@ $(HEADER) "=====> Signing and tagging the release..."
	@ if [ $$(git status -s -uall | wc -l) != 1 ]; then $(ERROR) "Git workspace must be clean."; exit 1; fi;

	@ $(WHITE) "This release will be tagged as: $(NEXT_VERSION)"
	@ echo "---------------------------------------------------------------------"
	@ read -p "Press any key to continue, or press Control+C to cancel. " x;

	@ echo " "
	@ chag update $(NEXT_VERSION)
	@ echo " "

	@ $(HEADER) "These are the contents of the CHANGELOG for this release. Are these correct?"
	@ $(WHITE) "---------------------------------------------------------------------"
	@ chag contents
	@ $(WHITE) "---------------------------------------------------------------------"
	@ echo "Are these release notes correct? If not, cancel and update CHANGELOG.md."
	@ read -p "Press any key to continue, or press Control+C to cancel. " x;

	@ echo " "

	git add .
	git commit -a -m "relprep: Preparing the $(NEXT_VERSION) release." --no-verify
	chag tag --sign

# goplicate-end:git

#-------------------------------------------------------------------------------
# Specific to this project.

.PHONY: clean
## clean: [clean]* Run standard cleanup tasks.
clean: clean-ds

.PHONY: lint
## lint: [lint]* Run linting tasks.
lint: pre-commit

.PHONY: test
## test: [test]* Runs ALL tests.
test:

#-------------------------------------------------------------------------------
# Installation

.PHONY: install-tools-go
## install-tools-go: [tools]* Install/upgrade the required Go packages.
install-tools-go:
	@ $(HEADER) "=====> Installing Go packages..."
	$(GO) install github.com/antham/gommit@latest
	$(GO) install github.com/google/yamlfmt/cmd/yamlfmt@latest
	$(GO) install github.com/pelletier/go-toml/v2/cmd/tomljson@latest

.PHONY: install-tools-mac
## install-tools-mac: [tools]* Install/upgrade the required tools for macOS, including Go packages.
install-tools-mac: install-tools-go
	@ $(HEADER) "=====> Installing required packages for macOS (Homebrew)..."
	brew update && brew install $(HOMEBREW_PACKAGES) && brew upgrade $(HOMEBREW_PACKAGES)
	curl -sSLf https://raw.githubusercontent.com/mtdowling/chag/master/install.sh | sudo bash

	@ $(BORDER) "To update to the latest versions, run:" "    brew update && brew upgrade"
