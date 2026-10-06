# gh-prune

Delete merged branches from GitHub repositories. A CLI tool to clean up merged branches automatically.

## Features

- Lists branches merged into the default branch
- Filters by age (only delete branches merged N+ days ago)
- Protects main branches (main, master, develop, release/*)
- Dry-run mode by default — preview before deleting
- Uses GitHub API directly, no gh CLI required

## Install

```bash
# From source
go install github.com/sablelight/gh-prune@latest

# Or download binary from releases
```

## Usage

```bash
# Dry run (default) - preview what would be deleted
gh-prune -o myorg -r myrepo -t $GH_TOKEN --dry-run

# Actually delete branches merged more than 7 days ago
gh-prune -o myorg -r myrepo -t $GH_TOKEN --dry-run=false

# Custom age threshold (30 days)
gh-prune -o myorg -r myrepo --days 30 --dry-run=false

# Custom protected branches
gh-prune -o myorg -r myrepo --protect main,master,develop,staging --dry-run=false
```

## Flags

| Flag | Short | Description | Default |
|---|---|---|---|
| `--owner` | `-o` | Repository owner (required) | |
| `--repo` | `-r` | Repository name (required) | |
| `--token` | `-t` | GitHub token | `GH_TOKEN` env |
| `--dry-run` | | Preview only, don't delete | `true` |
| `--days` | `-d` | Min days since merge | `7` |
| `--protect` | `-p` | Branches to never delete | `main,master,develop,release` |

## Authentication

Provide a GitHub token with `repo` scope:

```bash
export GH_TOKEN=ghp_xxxxxxxxxxxx
gh-prune -o myorg -r myrepo
```

Or pass directly:

```bash
gh-prune -o myorg -r myrepo -t ghp_xxxxxxxxxxxx
```

## Docker

```bash
docker run --rm -e GH_TOKEN=$GH_TOKEN \
  ghcr.io/sablelight/gh-prune -o myorg -r myrepo --dry-run=false
```

## How it works

1. Fetches the repository's default branch
2. Lists all branches (paginated)
3. For each branch, checks if there's a merged PR from that branch
4. Filters by age and protection rules
5. Deletes the branch ref via GitHub API

## License

MIT — see [LICENSE](LICENSE).