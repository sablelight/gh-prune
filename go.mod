module github.com/sablelight/gh-prune

go 1.23

require (
	github.com/google/go-github/v68 v68.0.0
	github.com/spf13/cobra v1.8.0
	github.com/sablelight/telegram-bot/config v0.0.0
)

replace (
	github.com/sablelight/telegram-bot/config => ../telegram-bot/config
)