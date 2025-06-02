
CODE_CHECK_TOOLS_SHELL="./scripts/code_check_tools.sh"
COPYRIGHT_SHELL="./scripts/copyright/update-copyright.sh"

$(shell chmod +x  ${CODE_CHECK_TOOLS_SHELL} ${COPYRIGHT_SHELL})

.PHONY: lint
lint:
	@${CODE_CHECK_TOOLS_SHELL} lint "$(mod)"
	@echo "lint check finished"

.PHONY: fix
fix: $(LINTER)
	@${CODE_CHECK_TOOLS_SHELL} fix "$(mod)"
	@echo "lint fix finished"

.PHONY: test
test:
	@${CODE_CHECK_TOOLS_SHELL} test "$(mod)"
	@echo "test check finished"

.PHONY: test-coverage
test-coverage:
	@${CODE_CHECK_TOOLS_SHELL} test_coverage "$(mod)"
	@echo "go test with coverage finished"

.PHONY: tidy
tidy:
	@${CODE_CHECK_TOOLS_SHELL} tidy
	@echo "tidy finished"

.PHONY: copyright
copyright:
	@${COPYRIGHT_SHELL}
	@echo "add copyright finish"

.PHONY: chglog
chglog:
	@git-chglog