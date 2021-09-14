package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2"
)

// CopyBox is to copy the contents of packr.Box to dst, except for the ignores list.
func CopyBox(dst io.Writer, box *packr.Box, ignores []string, format ...string) (int64, error) {
	files := box.List()
	var count int64
	for _, f := range files {
		isBreak := false
		for _, ig := range ignores {
			if f == ig {
				isBreak = true
				break
			}
		}
		if isBreak {
			continue
		}
		i, err := CopyBoxFile(dst, box, f, format...)
		if err != nil {
			return count, err
		}
		count += i
	}

	return count, nil
}

// CopyBoxFile copies the content of the specified name in packr.Box to dst.
func CopyBoxFile(dst io.Writer, box *packr.Box, name string, format ...string) (int64, error) {
	bs, err := box.Find(name)
	if err != nil {
		return 0, err
	}

	if 1 == len(format) {
		_, fullname := filepath.Split(name)
		var extension = filepath.Ext(fullname)
		var filename = fullname[0 : len(fullname)-len(extension)]
		bs = []byte(fmt.Sprintf(format[0], fullname, filename, string(bs)))
	} else if 1 < len(format) {
		return 0, errors.New("invaild format")
	}

	count, err := dst.Write(bs)

	return int64(count), err
}

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
	f.Sync()
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
	f.Sync()
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

// CopyDir copies from src dir to dst until either EOF is reached
// on src or an error occurs
func CopyDir(dst io.Writer, src string, format string) (int64, error) {
	fds, err := ioutil.ReadDir(src)
	if err != nil {
		return 0, err
	}

	var count int64
	for _, fd := range fds {
		fdPath := path.Join(src, fd.Name())

		var i int64
		var err error
		if fd.IsDir() {
			i, err = CopyDir(dst, fdPath, format)
		} else {
			i, err = CopyFile(dst, fdPath, format)
		}

		if err != nil {
			return count, err
		}

		count += i
	}

	return count, nil
}

// CopyFile copies from src dir to dst until either EOF is reached
// on src or an error occurs
func CopyFile(dst io.Writer, src string, format string) (int64, error) {
	if "" == format {
		srcfd, err := os.Open(src)
		if err != nil {
			return 0, err
		}
		defer srcfd.Close()

		return io.Copy(dst, srcfd)
	}

	bs, err := ioutil.ReadFile(src)
	if err != nil {
		return 0, err
	}

	_, filename := filepath.Split(src)
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	bs = []byte(fmt.Sprintf(format, filename, name, string(bs)))

	count, err := dst.Write(bs)

	return int64(count), err
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
