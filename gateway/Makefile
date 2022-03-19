BUF_VERSION:=1.1.0
MIGRATION_DIR=scripts/db/migrations

generate:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) generate

lint:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) lint
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) breaking --against 'https://github.com/overridesh/checkout-challenger.git#branch=master'
