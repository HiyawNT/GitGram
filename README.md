# GitGram  
**Telegram bot → GitHub push notifications**  
Get instant notifications in Telegram whenever someone pushes code to your favorite repositories.


## Features

- Subscribe to any public GitHub repository with `/subscribe owner/repo`
- Receive clean notifications for every new commit/push event
- Support for multiple subscriptions per user
- Unsubscribe anytime with `/unsubscribe`
- List all your subscriptions with `/list_subscriptions`
- Persistent storage using SQLite
- Built with **Go + Gin + go-telegram-bot-api + GORM**

## Example Notification
📢 New commit in hiyawNT/awesome-project
👤 By: hiyaw
📝 Fix login bug & improve rate limiting
🔗 https://github.com/hiyawNT/awesome-project/commit/abc123...


## Quick Start

### 1. Prerequisites

- Go 1.21+
- Telegram Bot Token (get it from [@BotFather](https://t.me/botfather))
- (Optional) A public server or tunnel (ngrok, fly.io, railway, render.com, etc.)

### 2. Setup

```bash
# Clone the project
git clone https://github.com/HiyawNT/GitGram.git
cd GitGram

# Copy example env file
touch  .env

# Edit .env (VERY IMPORTANT)
# ────────────────────────────────────────
# TELEGRAM_TOKEN=your_bot_token_here
# PORT=8080                 # change if needed
# ────────────────────────────────────────
```
### 3. Run locally 

```bash 

go mod tidy
go run .

```

### 4. Expose your server publicly (for local run)

use `ngrok` (quickest for testing):

```bash
ngrok http 8080

```
Copy the https://xxxx.ngrok.io URL → you'll need it for the GitHub webhook.

### 5. Add webhook in GitHub (per repository)

1. Go to repository → Settings → Webhooks → Add webhook
2. Payload URL: https://xxxx.ngrok.io/webhook (or your real domain)
3. Events: Select Just the push event
4. Active: checked
5. Click Add webhook

**Done**: Now push something → you should receive a message in Telegram.

### Bot Commands

**`/start`** : Welcome message

**`/subscribe owner/repo`**  : Subscribe to a GitHub repository

**`/unsubscribe owner/repo`** : Unsubscribe from a repository

**`/list_subscriptions`** : Show all your active subscriptions
**`/help`** : Show available commands

### Project Structure
```
GitGram/
├── main.go
├── handlers/
│   ├── github.go          ← GitHub webhook handler
│   └── telegram_handlers.go ← Telegram commands
├── models/
│   ├── github_payload.go
│   └── subscription.go
├── services/
│   └── telegram.go
├── storage/
│   └── subscription.go    ← GORM + SQLite logic
├── utils/
│   └── config.go
└── gitgram.db             ← created automatically

```

