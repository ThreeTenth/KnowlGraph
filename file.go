package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// AppendFile append src file to dst
func AppendFile(src, dst string) error {
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	_, err = f.Write(bs)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// MergeFiles merge multiple files in the specified directory to the specified file
func MergeFiles(src, dst string, before, after string, format, separator string, ignores ...string) error {
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = f.WriteString(before)
	if err != nil {
		return err
	}
	err = mergeFiles(src, f, format, separator, ignores...)
	_, err = f.WriteString(after)
	if err != nil {
		return err
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func mergeFiles(src string, dst *os.File, format string, separator string, ignores ...string) error {
	st, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return appendFile(src, dst, format, separator)
	}

	jsFds, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fd := range jsFds {
		fdPath := path.Join(src, fd.Name())
		ok := true

		for _, ig := range ignores {
			if fdPath == ig || strings.Contains(fdPath, ig) {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		if fd.IsDir() {
			if err = mergeFiles(fdPath, dst, format, separator, ignores...); err != nil {
				return err
			}
		} else if err = appendFile(fdPath, dst, format, separator); err != nil {
			return err
		}
	}

	return nil
}

func appendFile(src string, dst *os.File, format string, separator string) error {
	fd, err := os.Stat(src)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if format != "" {
		var filename = fd.Name()
		var extension = filepath.Ext(filename)
		var name = filename[0 : len(filename)-len(extension)]
		bs = []byte(fmt.Sprintf(format, fd.Name(), name, string(bs)))
	}
	if separator != "" {
		bs = append(bs, separator...)
	}
	_, err = dst.Write(bs)

	return err
}

// CpDir copies a whole directory recursively
func CpDir(src, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CpDir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = CpFile(srcfp, dstfp); err != nil {
				return err
			}
		}
	}

	return nil
}

// CpFile copies a single file from src to dst
func CpFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
