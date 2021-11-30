init:
	@echo "\
	GROUPWARE_ID=\n\
	GROUPWARE_PW=\n\
	SLACK_BOT_TOKEN=\n\
	" >> .env

# crawl groupware -job declined_payments
# crawl hackernews --channel hacker-news --point_threshold 100
# crawl quasarzone --channel quasarzone
# crawl book --channel gos16052
# crawl eomisae --channel gos16052 --target raffle
up:
	@COMMAND="$(cmd)" docker-compose up --build app chrome

# bin/crawl-data-slack crawl groupware -job declined_payments
# bin/crawl-data-slack crawl hackernews --channel hacker-news --point_threshold 100
# bin/crawl-data-slack crawl quasarzone --channel quasarzone
# bin/crawl-data-slack crawl book --channel gos16052
# bin/crawl-data-slack crawl eomisae --channel gos16052 --target raffle
shell:
	@docker-compose build app
	@docker-compose run --name crawl-data-slack-shell --rm app bash

docker-upload:
	@docker build .
	@docker tag crawl-data-slack_app gos16052/crawl-data-slack:$(version)
	@docker push gos16052/crawl-data-slack:$(version)

helm-upgrade: docker-upload
	@helm upgrade grslack helm/crawl-data-slack -f helm/values/values.yaml

helm-template:
	@helm template grslack helm/crawl-data-slack -f helm/values/values.yaml
