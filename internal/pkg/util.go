package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// exists returns whether the given file or directory exists
func ExistsFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ExePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func ExeDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// Check the slice contains the element
// s Slice or arr
// e element
func Contains(s interface{}, e interface{}) bool {
	arr := reflect.ValueOf(s)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == e {
			return true
		}
	}
	return false
}
