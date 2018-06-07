package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst. attempting to preserve permissions.
// If dst flie path directory not exit, will make it.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	// make dir
	fi, e := os.Stat(filepath.Dir(src))
	if e == nil {
		e = MakeDirWithMode(filepath.Dir(dst), fi.Mode())
	} else {
		e = MakeDir(filepath.Dir(dst))
	}
	if e != nil {
		return e
	}

	// copy
	err = copyFileContents(src, dst)

	// attempting to preserve permissions
	if err == nil {
		os.Chmod(dst, sfi.Mode())
	}

	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// CopyDir Recursively copies a directory tree, and preserve permissions.
// source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return fmt.Errorf("Destination already exists")
	}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {
		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				return
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				return
			}
		}
	}
	return
}

// MakeDir make a directory.
// if dest exist, just return nil.
func MakeDir(dest string) error {
	return MakeDirWithMode(dest, os.ModePerm)
}

// MakeDirWithMode make a directory.
// if dest exist, just return nil.
func MakeDirWithMode(dest string, permissions os.FileMode) error {
	_, err := os.Stat(dest)
	if err == nil { // 文件或文件夹已存在
		return nil
	} else if os.IsNotExist(err) { // 文件或文件夹不存在
		return os.MkdirAll(dest, permissions)
	} else {
		return err
	}
}
