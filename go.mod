module github.com/sablelight/gh-prune

go 1.23

require (
	github.com/google/go-github/v68 v68.0.0
	github.com/sablelight/telegram-bot/config v0.0.0
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

replace github.com/sablelight/telegram-bot/config => ../telegram-bot/config
