export SQLITE_VSS_EXT_PATH=tmplib

default: test

ci: depsdev test

test:
	go test ./... -coverprofile=coverage.out -covermode=count

lint:
	golangci-lint run ./...

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	gh ext install k1LoW/gh-setup --force
	rm -rf ${SQLITE_VSS_EXT_PATH}
	mkdir ${SQLITE_VSS_EXT_PATH}
	gh setup --repo github.com/asg017/sqlite-vss --bin-dir ${SQLITE_VSS_EXT_PATH} --match sqlite-vss-v.*-loadable --skip-content-type-check

prerelease_for_tagpr: depsdev
	go mod tidy
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag

.PHONY: default test
