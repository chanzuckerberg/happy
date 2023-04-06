package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DirProcessor interface {
	Tarzip(src string, f *os.File) error
}

type LocalProcessor struct{}

func NewLocalProcessor() *LocalProcessor {
	return &LocalProcessor{}
}

func (s *LocalProcessor) Tarzip(src string, f *os.File) error {
	logrus.Debugf("Tarzipping file (%s) ...", f.Name())

	if _, err := os.Stat(src); err != nil {
		return errors.Errorf("fail to tar file: %v", err)
	}
	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		dir := strings.Split(filepath.Dir(file), "/")[0]
		logrus.Debugf("Processing file %s (%s) ...", fi.Name(), file)
		if err != nil {
			return errors.Wrapf(err, "Unable to walk the file path %s", file)
		}

		if !fi.Mode().IsRegular() {
			logrus.Debugf("skipping file (%s) ...", fi.Name())
			return nil
		}

		if fi.Name() == ".DS_Store" || fi.Name() == ".terraform" || fi.Name() == ".git" || fi.Name() == ".terraform.lock.hcl" || filepath.Ext(fi.Name()) == ".tar.gz" || dir == ".terraform" || dir == ".git" {
			logrus.Debugf("Skipping file (%s) ...", fi.Name())
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.Wrapf(err, "Failed to get the file info header %s", fi.Name())
		}

		header.Name = file

		if err := tw.WriteHeader(header); err != nil {
			return errors.Wrapf(err, "Failed to write the file header %s", header.Name)
		}

		f, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "Cannot open file %s", file)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.Wrap(err, "Cannot copy file")
		}

		f.Close()

		return nil
	})
}
