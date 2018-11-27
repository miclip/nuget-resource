package nugetresource

import (
	"compress/flate"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"

	"github.com/mitchellh/colorstring"
)

func Fatal(doing string, err error) {
	Sayf(colorstring.Color("[red]error %s: %s \n "), doing, err)
	os.Exit(1)
}

func Sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}

// MakeZip creates a zip archive file
func MakeZip(zipName string, files []string) error {

	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	err := z.Archive(files, zipName)
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveZip ...
func UnarchiveZip(zipName string, destination string) error {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: false,
	}
	err := z.Unarchive(zipName, destination)
	if err != nil {
		return err
	}
	return nil
}

func ChmodAllFiles(directory string, mode os.FileMode) {
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			err = os.Chmod(path, mode)
			if err != nil {
				Sayf("chmod failed %v", err)
			}
			return nil
		})
	if err != nil {
		Sayf("chmod failed %v", err)
	}
}
