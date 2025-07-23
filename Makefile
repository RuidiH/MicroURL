.PHONY: run-local run-prod clean del-local del-prod

run-local:
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

run-prod:
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

del-local:
	@make clean
	@cd terraform && terraform init \
	-backend-config=envs/local/backend.conf \
	-reconfigure && \
	terraform workspace select local && \
	terraform destroy \
	-var-file=terraform.tfvars.local
	-auto-approve && \
	terraform workspace select default && \
	terraform workspace delete local

del-prod:
	@make clean
	@cd terraform && terraform init \
	-backend-config=envs/prod/backend.conf \
	-reconfigure && \
	terraform workspace select prod && \
	terraform destroy \
	-auto-approve && \
	terraform workspace select default && \
	terraform workspace delete prod 