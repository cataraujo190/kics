package source

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/checkmarxDev/ice/internal/logger"
	"github.com/pkg/errors"
)

const extensionTerraform = ".tf"

type FileSystemSourceProvider struct {
	Path string
}

var ErrNotSupportedFile = errors.New("invalid file format")

func (s *FileSystemSourceProvider) GetSources(ctx context.Context, scanID string, sink Sink) error {
	fileInfo, err := os.Stat(s.Path)
	if err != nil {
		return errors.Wrap(err, "failed to open path")
	}

	if !fileInfo.IsDir() {
		if filepath.Ext(s.Path) != extensionTerraform {
			return ErrNotSupportedFile
		}

		c, errOpenFile := os.Open(s.Path)
		if errOpenFile != nil {
			return errors.Wrap(errOpenFile, "failed to open path")
		}

		return sink(ctx, fileInfo.Name(), c)
	}

	err = filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != extensionTerraform {
			return nil
		}

		c, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}

		err = sink(ctx, strings.ReplaceAll(path, "\\", "/"), c)
		if err != nil {
			logger.GetLoggerWithFieldsFromContext(ctx).
				Err(err).
				Msgf("Filesystem terraform files provider couldn't parse file, file=%s", info.Name())
		}

		return nil
	})

	return errors.Wrap(err, "failed to walk directory")
}