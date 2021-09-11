crawl-data-slack
=======
docker-hub [gos16052/crawl-data-slack](https://hub.docker.com/r/gos16052/crawl-data-slack)

# Setup
## Requirements
- mysql (migrations are run automatically)

# Feed [hackernews](https://news.ycombinator.com/news)

![img](https://user-images.githubusercontent.com/13393411/132940346-4753e779-a5bb-434a-9a75-f0d3d3df0254.png)


1. Create new app (https://api.slack.com/apps)
2. Go to your app - OAuth & Permissions
3. Add thescopes into bot token scopes
(`channels.read, chat:write, groups:read, im:read, mpim:read, users:read`)
It uses 3 apis ([conversations.list](https://api.slack.com/methods/conversations.list), [chat.postMessage](https://api.slack.com/methods/chat.postMessage), [users.list](https://api.slack.com/methods/users.list))
![img](https://user-images.githubusercontent.com/13393411/132941304-5388ddfd-85eb-4fa2-97e9-28c51e6463e4.png)
4. Install to workspace and copy your bot user token

![image](https://user-images.githubusercontent.com/13393411/132941333-21a2e9f3-8c48-43b2-bbe2-c640aa33e506.png)

## Run the crawler

1. Pull the repository to run with chromedp
2. Creata a `.env.` file and write the environments
```
# Fill ID, PW, HOST, PORT, DATABASE, SLACK_BOT_TOKEN
MYSQL_CONN=<ID>:<PW>@tcp(<HOST>:<PORT>)/<DATABASE>?charset=utf8&parseTime=True
SLACK_BOT_TOKEN=<SLACK_BOT_TOKEN>
```

3. Run with the command `make up cmd="crawl\ hackernews\ --chanenl\ <YOUR_CHANNEL>`
