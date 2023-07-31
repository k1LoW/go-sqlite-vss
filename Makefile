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

testdata:
	mkdir -p testdata
	curl https://dl.fbaipublicfiles.com/fasttext/vectors-crawl/cc.en.300.vec.gz -o testdata/cc.en.300.vec.gz
	cd testdata && gzip -d cc.en.300.vec.gz
	head -1001  testdata/cc.en.300.vec | tail +2 > testdata/test.vec
	rm testdata/cc.en.300.vec.gz

.PHONY: default test testdata
