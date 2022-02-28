##@ General

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

init: ## init creates .env file to inject environment variables
	@echo "\
	GROUPWARE_ID=\n\
	GROUPWARE_PW=\n\

	EOMISAE_ID=\n\
	EOMISAE_PW=\n\

	SLACK_BOT_TOKEN=\n\
	" >> .env

# Examples
# make up cmd="crawl groupware -job declined_payments"
# make up cmd="crawl hackernews --channel CHANNEL --point_threshold 100"
# make up cmd="crawl quasarzone --channel CHANNEL"
# make up cmd="crawl eomisae --channel CHANNEL --target raffle"
# make up cmd="crawl financial-report --channel CHANNEL"
# make up cmd="crawl gitpublic --channel CHANNEL"
# make up cmd="crawl ipo --channel CHANNEL"
# make up cmd="crawl gitpublic --channel CHANNEL"
# make up cmd="crawl spinnaker --channel CHANNEL --token GITHUB_TOKEN --host HOST"
# make up cmd="crawl wanted --channel CHANNEL --query \"SEARCH_STRING\""
# make up cmd="crawl techcrunch --channel CHANNEL"
# make up cmd="crawl guardian --channel CHANNEL --url 'https://www.theguardian.com/world/live/2022/feb/27/russia-ukraine-latest-news-missile-strikes-on-oil-facilities-reported-as-some-russian-banks-cut-off-from-swift-system-live?ilterKeyEvents=true'"
up: ## Run the application `make up cmd="crawl financial-report --channel my_channel"`, open the Makefile to see more examples.
	@docker-compose build app
	@COMMAND="$(cmd)" docker-compose up app

# crawl-data-slack crawl groupware -job declined_payments
# crawl-data-slack crawl hackernews --channel hacker-news --point_threshold 100
# crawl-data-slack crawl quasarzone --channel quasarzone
# crawl-data-slack crawl book --channel gos16052
# crawl-data-slack crawl eomisae --channel gos16052 --target raffle
# crawl-data-slack crawl ipo --channel gos16052
shell: ## Run the shell
	@docker-compose build app
	@docker-compose run --name crawl-data-slack-shell --rm app bash

run: ## Run the cmd
	@docker-compose build app
	@docker-compose run --name crawl-data-slack-run --rm app $(cmd)

docker-upload: ## Upload the image to the docker registry `make docker-upload version=0.x.y`
	@docker-compose build
	@docker tag crawl-data-slack_app gos16052/crawl-data-slack:$(version)
	@docker push gos16052/crawl-data-slack:$(version)

# eks-rw
# kubectx dev
# kubens grslack
helm-upgrade: docker-upload
	@helm upgrade grslack helm/crawl-data-slack -f helm/values/values.yaml

helm-template:
	@helm template grslack helm/crawl-data-slack -f helm/values/values.yaml
