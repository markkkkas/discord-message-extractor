# Discord Message Extractor

Simple HTTP service, for extracting messages from Discord server.

### Usage
```bash
> go build main.go
> ./main -token DISCORD_BOT_TOKEN
```

```bash
> curl http://localhost:8080?channelId=CHANNEL_ID&limit=10 
```