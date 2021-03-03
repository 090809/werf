package true_git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/werf/logboek"
)

type SyncSourceWorktreeWithServiceWorktreeBranchOptions struct {
	OnlyStagedChanges bool
}

func SyncSourceWorktreeWithServiceWorktreeBranch(ctx context.Context, gitDir, sourceWorkTreeDir, workTreeCacheDir, commit string, opts SyncSourceWorktreeWithServiceWorktreeBranchOptions) (string, error) {
	var resultCommit string
	if err := withWorkTreeCacheLock(ctx, workTreeCacheDir, func() error {
		var err error
		if gitDir, err = filepath.Abs(gitDir); err != nil {
			return fmt.Errorf("bad git dir %s: %s", gitDir, err)
		}

		if workTreeCacheDir, err = filepath.Abs(workTreeCacheDir); err != nil {
			return fmt.Errorf("bad work tree cache dir %s: %s", workTreeCacheDir, err)
		}

		destinationWorkTreeDir, err := prepareWorkTree(ctx, gitDir, workTreeCacheDir, commit, true)
		if err != nil {
			return fmt.Errorf("unable to prepare worktree for commit %v: %s", commit, err)
		}

		currentCommitPath := filepath.Join(workTreeCacheDir, "current_commit")
		if err := os.RemoveAll(currentCommitPath); err != nil {
			return fmt.Errorf("unable to remove %s: %s", currentCommitPath, err)
		}

		devBranchName := fmt.Sprintf("werf-dev-%s", commit)
		devCommitWithStagedChanges, err := syncWorktreeWithServiceWorktreeBranch(ctx, sourceWorkTreeDir, destinationWorkTreeDir, commit, devBranchName, true)
		if err != nil {
			return fmt.Errorf("unable to sync staged changes: %s", err)
		}

		if opts.OnlyStagedChanges {
			resultCommit = devCommitWithStagedChanges
		} else {
			devBranchName = fmt.Sprintf("werf-dev-%s-%s", commit, devCommitWithStagedChanges)
			devCommitWithTrackedChanges, err := syncWorktreeWithServiceWorktreeBranch(ctx, sourceWorkTreeDir, destinationWorkTreeDir, devCommitWithStagedChanges, devBranchName, false)
			if err != nil {
				return fmt.Errorf("unable to sync tracked changes: %s", err)
			}

			resultCommit = devCommitWithTrackedChanges
		}

		return nil
	}); err != nil {
		return "", err
	}

	return resultCommit, nil
}

func syncWorktreeWithServiceWorktreeBranch(ctx context.Context, sourceWorkTreeDir, destinationWorkTreeDir, commit string, devBranchName string, onlyStagedChanges bool) (string, error) {
	var isDevBranchExist bool
	if output, err := runGitCmd(ctx, []string{"branch", "--list", devBranchName}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
		return "", err
	} else {
		isDevBranchExist = output.Len() != 0
	}

	var devHeadCommit string
	if isDevBranchExist {
		if _, err := runGitCmd(ctx, []string{"checkout", devBranchName}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
			return "", err
		}

		if output, err := runGitCmd(ctx, []string{"rev-parse", devBranchName}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
			return "", err
		} else {
			devHeadCommit = strings.TrimSpace(output.String())
		}
	} else {
		if _, err := runGitCmd(ctx, []string{"checkout", "-b", devBranchName, commit}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
			return "", err
		}

		devHeadCommit = commit
	}

	gitDiffArgs := []string{
		"-c", "diff.renames=false",
		"-c", "core.quotePath=false",
		"diff",
		"--full-index",
		"--binary",
	}
	if onlyStagedChanges {
		gitDiffArgs = append(gitDiffArgs, "--cached")
	}
	gitDiffArgs = append(gitDiffArgs, devHeadCommit)

	var resCommit string
	if diffOutput, err := runGitCmd(ctx, gitDiffArgs, sourceWorkTreeDir, runGitCmdOptions{}); err != nil {
		return "", err
	} else if len(diffOutput.Bytes()) == 0 {
		if debug() {
			logboek.Context(ctx).Debug().LogLn("[DEBUG] Nothing to sync")
		}
		resCommit = devHeadCommit
	} else {
		if debug() {
			filesType := "tracked"
			if onlyStagedChanges {
				filesType = "staged"
			}
			logboek.Context(ctx).Debug().LogF("[DEBUG] Syncing %s files ...\n", filesType)
		}

		if _, err := runGitCmd(ctx, []string{"apply", "--binary", "--index"}, destinationWorkTreeDir, runGitCmdOptions{stdin: diffOutput}); err != nil {
			return "", err
		}

		gitArgs := []string{"-c", "user.email=werf@werf.io", "-c", "user.name=werf", "commit", "-m", time.Now().String()}
		if _, err := runGitCmd(ctx, gitArgs, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
			return "", err
		}

		if output, err := runGitCmd(ctx, []string{"rev-parse", devBranchName}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
			return "", err
		} else {
			newDevCommit := strings.TrimSpace(output.String())
			resCommit = newDevCommit
		}
	}

	if _, err := runGitCmd(ctx, []string{"checkout", "--force", "--detach", resCommit}, destinationWorkTreeDir, runGitCmdOptions{}); err != nil {
		return "", err
	}

	return resCommit, nil
}

type runGitCmdOptions struct {
	stdin io.Reader
}

func runGitCmd(ctx context.Context, args []string, dir string, opts runGitCmdOptions) (*bytes.Buffer, error) {
	allArgs := append(getCommonGitOptions(), args...)
	cmd := exec.Command("git", allArgs...)
	cmd.Dir = dir

	if opts.stdin != nil {
		cmd.Stdin = opts.stdin
	}

	output := SetCommandRecordingLiveOutput(ctx, cmd)

	err := cmd.Run()

	cmdWithArgs := strings.Join(append([]string{cmd.Path, "-C " + dir}, cmd.Args[1:]...), " ")
	if debug() {
		fmt.Printf("[DEBUG] %s\n%s\n", cmdWithArgs, output)
	}

	if err != nil {
		return nil, fmt.Errorf("git command %s failed: %s\n%s", cmdWithArgs, err, output)
	}

	return output, err
}

func debug() bool {
	return os.Getenv("WERF_DEBUG_TRUE_GIT") == "1"
}
