# GitGram - GitHub Push Notifications via Telegram

GitGram is a Telegram bot that delivers GitHub repository push notifications directly to your Telegram chats. Subscribe to repositories and receive instant updates when commits are pushed.

## Features

- 📬 **Subscribe to Repositories**: Get notifications for specific GitHub repositories
- 🔔 **Real-time Updates**: Receive push notifications instantly via Telegram
- 📋 **Subscription Management**: Easily manage your subscriptions
- 🔒 **Concurrent Safe**: Thread-safe storage with mutex protection
- 🚀 **Easy Setup**: Simple configuration with environment variables

## Prerequisites

- Go 1.21 or higher
- A Telegram Bot Token (from [@BotFather](https://t.me/botfather))
- A publicly accessible server to receive GitHub webhooks

## Installation

1. **Clone the repository**:

```bash
git clone <your-repo-url>
cd gitgram
```

2. **Install dependencies**:

```bash
go mod download
```

3. **Create a `.env` file**:

```bash
cp .env.example .env
```

4. **Add your Telegram Bot Token to `.env`**:

```
TELEGRAM_TOKEN=your_telegram_bot_token_here
PORT=8080
```

5. **Run the bot**:

```bash
go run main.go
```

## Usage

### Telegram Bot Commands

- `/start` - Welcome message and list of commands
- `/help` - Display help information
- `/subscribe <owner/repo>` - Subscribe to a repository
  - Example: `/subscribe octocat/Hello-World`
- `/unsubscribe <owner/repo>` - Unsubscribe from a repository
- `/list_subscriptions` - View your active subscriptions

### Setting Up GitHub Webhooks

To receive push notifications, you need to configure a webhook in your GitHub repository:

1. Go to your GitHub repository
2. Navigate to **Settings** → **Webhooks** → **Add webhook**
3. Configure the webhook:
   - **Payload URL**: `https://your-server.com/webhook`
   - **Content type**: `application/json`
   - **Which events**: Select "Just the push event"
   - **Active**: ✓ Checked
4. Click **Add webhook**

The bot will now send notifications to all subscribed users when commits are pushed to this repository.

## Project Structure

```
gitgram/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── .env                         # Environment variables (not in git)
├── .env.example                 # Example environment file
├── utils/
│   └── config.go               # Configuration loader
├── models/
│   └── github_payload.go       # GitHub webhook payload structs
├── storage/
│   ├── subscription.go         # Subscription storage logic
│   └── subscription_test.go    # Storage unit tests
├── services/
│   └── telegram.go             # Telegram bot service
└── handlers/
    └── github.go               # GitHub webhook handler
```

## API Endpoints

- `GET /` - Health check endpoint
- `POST /webhook` - GitHub webhook receiver

## Testing

Run the unit tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests for a specific package:

```bash
go test ./storage -v
```

## Production Considerations

### Persistent Storage

The current implementation uses **in-memory storage**, which means all subscriptions will be lost when the application restarts. For production use, implement persistent storage:

**Recommended Options**:

- **PostgreSQL**: Best for large-scale deployments
- **MongoDB**: Good for document-based storage
- **SQLite**: Simple file-based database for smaller deployments
- **BoltDB**: Embedded key-value store for Go

To implement persistent storage, modify `storage/subscription.go` and replace the in-memory maps with database operations.

### Deployment

**Using systemd** (Linux):

Create a service file `/etc/systemd/system/gitgram.service`:

```ini
[Unit]
Description=GitGram Telegram Bot
After=network.target

[Service]
Type=simple
User=gitgram
WorkingDirectory=/opt/gitgram
ExecStart=/opt/gitgram/gitgram
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable gitgram
sudo systemctl start gitgram
```

**Using Docker**:

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o gitgram .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gitgram .
COPY .env .
CMD ["./gitgram"]
```

Build and run:

```bash
docker build -t gitgram .
docker run -p 8080:8080 --env-file .env gitgram
```

### Security

- Keep your `.env` file secure and never commit it to version control
- Use HTTPS for your webhook endpoint
- Consider implementing webhook signature verification (GitHub sends `X-Hub-Signature-256` header)
- Rate limit webhook requests to prevent abuse
- Run the bot with minimal privileges

### Monitoring

Consider adding:

- Prometheus metrics for monitoring
- Structured logging (e.g., using `logrus` or `zap`)
- Error tracking (e.g., Sentry)
- Health check endpoints

## Troubleshooting

**Bot doesn't respond to commands**:

- Verify the `TELEGRAM_TOKEN` is correct
- Check that the bot is running (`ps aux | grep gitgram`)
- Look at the logs for errors

**Webhooks not working**:

- Ensure your server is publicly accessible
- Check GitHub webhook delivery logs in repository settings
- Verify the webhook URL is correct
- Check firewall settings

**No notifications received**:

- Verify you're subscribed to the repository: `/list_subscriptions`
- Check the GitHub webhook is configured correctly
- Look at server logs for incoming webhook requests

## Development

To contribute to GitGram:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

MIT License - feel free to use this project for any purpose.

## Support

For issues, questions, or suggestions, please open an issue on the GitHub repository.

---

Built with ❤️ using Go, Gin, and the Telegram Bot API
