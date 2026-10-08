package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sablelight/gh-prune/internal/config"
	"github.com/sablelight/gh-prune/internal/github"
	"github.com/spf13/cobra"
)

var (
	owner string
	repo  string
	token string
	dryRun bool
	daysOld int
	protect []string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gh-prune",
		Short: "Delete merged branches from GitHub repositories",
		Long: `gh-prune deletes branches that have been merged into the default branch.
		
It uses the GitHub API to find merged branches and optionally deletes them.
Use --dry-run to preview what would be deleted without making changes.`,
		RunE: runPrune,
	}

	rootCmd.Flags().StringVarP(&owner, "owner", "o", "", "Repository owner (required)")
	rootCmd.Flags().StringVarP(&repo, "repo", "r", "", "Repository name (required)")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "GitHub token (or set GH_TOKEN env)")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", true, "Preview deletions without deleting")
	rootCmd.Flags().IntVar(&daysOld, "days", 7, "Only delete branches merged more than N days ago")
	rootCmd.Flags().StringSliceVar(&protect, "protect", []string{"main", "master", "develop", "release"}, "Branch names to never delete")

	_ = rootCmd.MarkFlagRequired("owner")
	_ = rootCmd.MarkFlagRequired("repo")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runPrune(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	ghToken := token
	if ghToken == "" {
		ghToken = config.RegisterOption("github.token", "GitHub token", "").GetString()
	}
	if ghToken == "" {
		return fmt.Errorf("GitHub token required (--token or GH_TOKEN env)")
	}

	client := github.NewClient(ghToken)

	// Get default branch
	defaultBranch, err := client.GetDefaultBranch(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("get default branch: %w", err)
	}
	fmt.Printf("Default branch: %s\n", defaultBranch)

	// Get merged branches
	branches, err := client.GetMergedBranches(ctx, owner, repo, defaultBranch)
	if err != nil {
		return fmt.Errorf("get merged branches: %w", err)
	}

	// Filter
	cutoff := time.Now().AddDate(0, 0, -daysOld)
	var toDelete []github.Branch
	for _, b := range branches {
		if isProtected(b.Name, protect) || strings.HasPrefix(b.Name, "release/") {
			continue
		}
		if b.MergedAt != nil && b.MergedAt.Before(cutoff) {
			toDelete = append(toDelete, b)
		}
	}

	if len(toDelete) == 0 {
		fmt.Println("No branches to delete.")
		return nil
	}

	fmt.Printf("\nBranches to delete (%d):\n", len(toDelete))
	for _, b := range toDelete {
		merged := "unknown"
		if b.MergedAt != nil {
			merged = b.MergedAt.Format("2006-01-02")
		}
		fmt.Printf("  %s (merged %s, SHA: %.7s)\n", b.Name, merged, b.SHA)
	}

	if dryRun {
		fmt.Println("\n[DRY RUN] No branches deleted. Run with --dry-run=false to delete.")
		return nil
	}

	// Delete
	fmt.Println("\nDeleting...")
	for _, b := range toDelete {
		if err := client.DeleteBranch(ctx, owner, repo, b.Name); err != nil {
			fmt.Printf("  ✗ %s: %v\n", b.Name, err)
		} else {
			fmt.Printf("  ✓ %s\n", b.Name)
		}
	}
	fmt.Println("\nDone.")
	return nil
}

func isProtected(name string, protected []string) bool {
	for _, p := range protected {
		if name == p {
			return true
		}
	}
	return false
}