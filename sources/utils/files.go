package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(filepath string, url string, message string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create progress bar
	bar := progressbar.DefaultBytes(resp.ContentLength, message)

	// Write the body to file with progress bar
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return err
}

func UnpackTarGz(gzipStream io.Reader, dst string, skip int) error {
	unzippedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer unzippedStream.Close()

	tarReader := tar.NewReader(unzippedStream)

	for {

		// Read next header
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		// Remove top level directory from path and append destination directory
		target := filepath.Join(dst, strings.Join(strings.Split(header.Name, "/")[skip:], "/"))

		// Handle dir/file according to header type
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
}
