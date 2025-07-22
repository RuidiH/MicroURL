.PHONY: local prod clean

local:
	@make clean
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

prod:
	@make clean
	@cd terraform && \
	terraform init \
	-backend-config=envs/prod/backend.conf \
	-reconfigure
	@terraform workspace select prod 2>/dev/null || \
	terraform workspace new prod 
	@cd terraform && \
	terraform apply \
		-auto-approve

clean:
	@cd terraform && rm -rf .terraform .terraform.lock.hcl