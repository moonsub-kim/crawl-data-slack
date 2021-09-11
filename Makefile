init:
	@echo "\
	GROUPWARE_ID=\n\
	GROUPWARE_PW=\n\
	SLACK_BOT_TOKEN=\n\
	" >> .env

up:
	@COMMAND=$(cmd) docker-compose up --build app chrome

cmd="crawl hackernews"
hackernews: up
	@echo "crawl hackernews"
