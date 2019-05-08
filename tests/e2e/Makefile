# This file should be included into main Makefile !
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

run-e2e:
	make test-e2e-env-start
	make test-e2e-iam
	make test-e2e-venom
	make test-e2e-env-stop


# test-e2e-env-start starts the test environment with docker-compose:
# 1. It starts nginx service
# 2. It starts redis service
test-e2e-env-start: ## Start test environment
	$(call before_job,Starting test environment:)
	cd $(ROOT_DIR) && docker-compose up -d --build nginx redis
	$(call after_job,Test environment started!)

# Builds the IAM container with required build arguments
build-iam:
	cd $(ROOT_DIR) && docker-compose build --build-arg GITLAB_USERNAME --build-arg GITLAB_PASSWORD kiwi-iam

# test-e2e-iam simply just starts iam service
test-e2e-iam: ## Start iam
	$(call before_job,Starting kiwi-iam:)
	cd $(ROOT_DIR) && docker-compose up -d kiwi-iam  
	$(call after_job,kiwi-iam started!)

test-e2e-venom:
	$(call before_job,Starting venom tests:)
	cd $(ROOT_DIR) && docker-compose build venom
	$(eval VENOM_CONTAINER_NAME = venom-container-$(shell date +%s))
	cd $(ROOT_DIR) && docker-compose run -d --name $(VENOM_CONTAINER_NAME) venom tail -f /dev/null
	docker exec $(VENOM_CONTAINER_NAME) venom run --var-from-file variables.yml --parallel 5 --format=xml --output-dir="." --strict ; \
	let "EXIT_CODE=$$?" ; \
	cd $(ROOT_DIR) && docker exec  $(VENOM_CONTAINER_NAME) cat test_results.xml > test_results.xml ; \
	docker-compose rm -f -v venom &1>/dev/null || true ; \
	exit $$EXIT_CODE
	$(call after_job,Venom integration tests completed!)

test-e2e-env-stop: ## Stop test environment
	docker-compose down -v