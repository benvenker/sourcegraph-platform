package check

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/sourcegraph/sourcegraph/dev/sg/internal/usershell"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type Check[Args any] struct {
	Name        string
	Description string

	// Check must be implemented to execute the check. Should be run using RunCheck.
	Check ActionFunc[Args]
	// Fix can be implemented to fix issues with this check.
	Fix ActionFunc[Args]
	// Enabled can be implemented to indicate when this checker should be skipped.
	Enabled EnableFunc[Args]

	// checkErr preserves the state of the most recent check run.
	checkErr error
	checkRun bool
}

// RunCheck should be used to run a check and set its results onto the Check itself.
func (c *Check[Args]) RunCheck(ctx context.Context, cio IO, args Args) error {
	c.checkRun = true
	c.checkErr = c.Check(ctx, cio, args)
	return c.checkErr
}

// IsMet indicates if this check has been run, and if it has errored. RunCheck should be
// called to update state.
func (c *Check[Args]) IsMet() bool {
	return c.checkRun && c.checkErr == nil
}

// Category is a set of checks.
type Category[Args any] struct {
	Name        string
	Description string
	Checks      []*Check[Args]

	// DependsOn lists names of Categories that must be fulfilled before checks in this
	// category are run.
	DependsOn []string

	// Enabled can be implemented to indicate when this checker should be skipped.
	Enabled EnableFunc[Args]
}

// HasFixable indicates if this category has any fixable checks.
func (c *Category[Args]) HasFixable() bool {
	for _, c := range c.Checks {
		if c.Fix != nil {
			return true
		}
	}
	return false
}

type CheckFunc func(context.Context) error

func InPath(cmd string) CheckFunc {
	return func(ctx context.Context) error {
		hashCmd := fmt.Sprintf("hash %s 2>/dev/null", cmd)
		_, err := usershell.CombinedExec(ctx, hashCmd)
		if err != nil {
			return errors.Newf("executable %q not found in $PATH", cmd)
		}
		return nil
	}
}

func CommandExitCode(cmd string, exitCode int) CheckFunc {
	return func(ctx context.Context) error {
		cmd := usershell.Cmd(ctx, cmd)
		err := cmd.Run()
		var execErr *exec.ExitError
		if err != nil {
			if errors.As(err, &execErr) && execErr.ExitCode() != exitCode {
				return errors.Newf("command %q has wrong exit code, wanted %d but got %d", cmd, exitCode, execErr.ExitCode())
			}
			return err
		}
		return nil
	}
}

func CommandOutputContains(cmd, contains string) CheckFunc {
	return func(ctx context.Context) error {
		out, _ := usershell.CombinedExec(ctx, cmd)
		if !strings.Contains(string(out), contains) {
			return errors.Newf("command output of %q doesn't contain %q", cmd, contains)
		}
		return nil
	}
}

func FileContains(fileName, content string) func(context.Context) error {
	return func(context.Context) error {
		file, err := os.Open(fileName)
		if err != nil {
			return errors.Wrapf(err, "failed to check that %q contains %q", fileName, content)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, content) {
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return errors.Newf("file %q did not contain %q", fileName, content)
	}
}

// This ties the check to having the library installed with apt-get on Ubuntu,
// which against the principle of checking dependencies independently of their
// installation method. Given they're just there for comby and sqlite, the chances
// that someone needs to install them in a different way is fairly low, making this
// check acceptable for the time being.
func HasUbuntuLibrary(name string) func(context.Context) error {
	return func(ctx context.Context) error {
		_, err := usershell.CombinedExec(ctx, fmt.Sprintf("dpkg -s %s", name))
		return errors.Newf("dpkg: %w", err)
	}
}

func Version(cmdName, haveVersion, versionConstraint string) error {
	c, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		return err
	}

	version, err := semver.NewVersion(haveVersion)
	if err != nil {
		return errors.Newf("cannot decode version in %q: %w", haveVersion, err)
	}

	if !c.Check(version) {
		return errors.Newf("version %q from %q does not match constraint %q", haveVersion, cmdName, versionConstraint)
	}
	return nil
}
