package util

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"syscall"
)

func Chroot(rootfs string) error {
	if err := syscall.Chroot(rootfs); err != nil {
		return errors.Wrapf(err, "unable to chroot to %s", rootfs)
	}

	root := filepath.FromSlash("/")
	if err := os.Chdir(root); err != nil {
		return errors.Wrapf(err, "unable to chdir to %s", root)
	}

	return nil
}
