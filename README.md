# Discord Messages Extractor

Simple HTTP service, for extracting messages from Discord server.

### Usage
```bash
> go build dme.go
> ./dme -token DISCORD_BOT_TOKEN
```

```bash
> curl http://localhost:8080?channelId=CHANNEL_ID&limit=10 
```