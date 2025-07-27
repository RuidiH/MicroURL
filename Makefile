.PHONY: run-local run-prod clean del-local del-prod build-lambdas build-create build-lookup

build-create:
	@mkdir -p src/create_url/bin
	@cd src/create_url && go mod tidy && \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/bootstrap main.go && \
	rm -f bin/create.zip && \
	cd bin && zip -j create.zip bootstrap

build-lookup:
	@mkdir -p src/lookup_redirect/bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o src/lookup_redirect/bin/bootstrap src/lookup_redirect/main.go
	cd src/lookup_redirect/bin && zip -j lookup.zip bootstrap

build-lambdas: build-create
# build-lambdas: build-create build-lookup

clean:
	@cd terraform && rm -rf .terraform .terraform.lock.hcl

run-local: clean build-lambdas
	@cd terraform && \
		terraform init \
		-backend-config=envs/local/backend.conf \
		-reconfigure
	@terraform workspace select local 2>/dev/null || \
		terraform workspace new local
	@cd terraform && \
	terraform apply \
		-var-file=terraform.tfvars.local \
		-auto-approve

run-prod: clean build-lambdas
	@cd terraform && \
	terraform init \
	-backend-config=envs/prod/backend.conf \
	-reconfigure
	@terraform workspace select prod 2>/dev/null || \
	terraform workspace new prod 
	@cd terraform && \
	terraform apply \
		-auto-approve

del-local: clean
	@cd terraform && terraform init \
	-backend-config=envs/local/backend.conf \
	-reconfigure && \
	terraform workspace select local && \
	terraform destroy \
	-var-file=terraform.tfvars.local
	-auto-approve && \
	terraform workspace select default && \
	terraform workspace delete local

del-prod: clean
	@cd terraform && terraform init \
	-backend-config=envs/prod/backend.conf \
	-reconfigure && \
	terraform workspace select prod && \
	terraform destroy \
	-auto-approve && \
	terraform workspace select default && \
	terraform workspace delete prod 