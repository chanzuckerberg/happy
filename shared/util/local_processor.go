package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ignoredEntries map[string]bool = map[string]bool{
	".DS_Store":           true,
	".terraform":          true,
	".git":                true,
	".terraform.lock.hcl": true,
}

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
		logrus.Debugf("Processing file %s (%s) ...", fi.Name(), file)
		if err != nil {
			return errors.Wrapf(err, "unable to walk the file path %s", file)
		}

		if _, ok := ignoredEntries[file]; ok {
			if fi.IsDir() {
				logrus.Debugf("Skipping folder (%s) ...", fi.Name())
				return filepath.SkipDir
			}
		}

		if !fi.Mode().IsRegular() {
			logrus.Debugf("Skipping file (%s) ...", fi.Name())
			return nil
		}

		if _, ok := ignoredEntries[fi.Name()]; ok {
			logrus.Debugf("Skipping file (%s) ...", fi.Name())
			return nil
		}

		if filepath.Ext(fi.Name()) == ".tar.gz" {
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
