##@ General

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

init: ## init creates .env file to inject environment variables
	@echo "\
	SLACK_BOT_TOKEN=\n\
	" >> .env

# Examples
# make up cmd="crawl --channel CHANNEL finance global-monitor"
# make up cmd="crawl --channel CHANNEL finance hankyung"
# make up cmd="crawl --channel CHANNEL finance mirae-asset"
# make up cmd="crawl --channel CHANNEL finance kcif"
# make up cmd="crawl --channel CHANNEL tech goldman-sachs --recent-days 3"
# make up cmd="crawl --channel CHANNEL tech naver-d2"
# make up cmd="crawl --channel CHANNEL tech hacker-news --point_threshold 100"
# make up cmd="crawl --channel CHANNEL tech quastor"
# make up cmd="crawl --channel CHANNEL tech delivery-hero"
# make up cmd="crawl --channel CHANNEL rss --name amazon-science --site https://www.amazon.science/index.rss --url-contains /latest-news/,/blog/ --recent-days 3"
# make up cmd="crawl --channel CHANNEL rss --name tech-blog-posts --site https://techblogposts.com/rss.xml --recent-days 1 --tech-blog-posts"
# make up cmd="crawl --channel CHANNEL lotte-cinema --date 2023-01-30"
# make up cmd="crawl --channel CHANNEL career wanted --query 'data analyst'"
# make up cmd="crawl --channel CHANNEL career wanted --query 'data scientist'"
# make up cmd="crawl --channel CHANNEL career wanted --query 'data engineer'"
# make up cmd="crawl --channel CHANNEL career wanted --query '데이터 사이언티스트'"
# make up cmd="crawl --channel CHANNEL career wanted --query 'brand design'"
# make up cmd="crawl --channel CHANNEL career wanted --query '브랜드 디자'"
# make up cmd="crawl --channel CHANNEL career wanted --query 'visual design'"
# make up cmd="crawl --channel CHANNEL career wanted --query '비주얼 디자'"
# make up cmd="crawl --channel CHANNEL career designer-job --query '그래픽 디자이너' --excepts '인턴'"
# make up cmd="crawl --channel CHANNEL career naver-career --query '데이터 사이언티스트'"
# make up cmd="slack archive --channel channel-name"

# make up cmd="crawl --channel CHANNEL finance ipo"
# make up cmd="crawl --channel CHANNEL confluent --job release"
# make up cmd="crawl --channel CHANNEL confluent --job status --channel kafka --keyword 'ap-northeast-1' --keyword 'Cloud\ Metrics' --keyword 'metrics API' --keyword 'ksqlDB' --keyword 'Confluent Cloud API'"
# make up cmd="crawl --channel CHANNEL rss --channel geeknews --name spotify --site 'https://engineering.atspotify.com/feed/' --recent-days 20"
up: ## Run the application `make up cmd="crawl --channel CHANNEL financial-report --channel my_channel"`, open the Makefile to see more examples.
	@docker-compose build app
	@COMMAND="$(cmd)" docker-compose up app

# crawl-data-slack crawl --channel CHANNEL groupware -job declined_payments
# crawl-data-slack crawl --channel CHANNEL hackernews --channel hacker-news --point_threshold 100
# crawl-data-slack crawl --channel CHANNEL quasarzone --channel quasarzone
# crawl-data-slack crawl --channel CHANNEL book --channel gos16052
# crawl-data-slack crawl --channel CHANNEL eomisae --channel gos16052 --target raffle
# crawl-data-slack crawl --channel CHANNEL ipo --channel gos16052
shell: ## Run the shell
	@docker-compose build app
	@docker-compose run --name crawl-data-slack-shell --rm app bash

run: ## Run the cmd
	@docker-compose build app
	@docker-compose run --name crawl-data-slack-run --rm app $(cmd)

docker-upload-m1: ## In apple silicon mac, upload the image to the docker registry `make docker-upload-m1 version=0.x.y`
	@docker buildx build --platform linux/amd64 --push -t gos16052/crawl-data-slack:$(version) .

docker-upload: ## Upload the image to the docker registry `make docker-upload version=0.x.y`
	@docker-compose build
	@docker tag crawl-data-slack_app gos16052/crawl-data-slack:$(version)
	@docker push gos16052/crawl-data-slack:$(version)

# eks-rw
# kubectx dev
# kubens grslack
helm-upgrade:
	@helm upgrade grslack helm/crawl-data-slack -f helm/values/values.yaml

helm-template:
	@helm template grslack helm/crawl-data-slack -f helm/values/values.yaml
