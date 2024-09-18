package validatecmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/acronis/go-raml/stacktrace"

	"github.com/acronis/go-cti/internal/app/cti"
	"github.com/acronis/go-cti/internal/pkg/command"
	_package "github.com/acronis/go-cti/pkg/package"
	"github.com/acronis/go-cti/pkg/pacman"
)

type cmd struct {
	opts    cti.Options
	targets []string
}

func New(opts cti.Options, targets []string) command.Command {
	return &cmd{
		opts:    opts,
		targets: targets,
	}
}

func (c *cmd) Execute(ctx context.Context) error {
	// workDir := filepath.Dir(c.targets[0])
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	p, err := pacman.New(filepath.Join(workDir, _package.IndexFileName))
	if err != nil {
		return err
	}
	errs := p.Validate()
	if errs != nil {
		for i := range errs {
			err := errs[i]
			if st, ok := stacktrace.Unwrap(err); ok {
				slog.Error(fmt.Sprintf("Tracebacks:\n%s", st.Sprint(stacktrace.WithEnsureDuplicates())))
			} else {
				slog.Error(err.Error())
			}
		}
		return errors.New("failed to validate the package")
	}
	slog.Info("No errors found")
	return nil
}
