package data

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-tools/record"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

var camSubDir = "cam"

func BuildArchive(basedir string, archiveName string, sliceSize int) error {
	dirItems, err := ioutil.ReadDir(basedir)
	if err != nil {
		return fmt.Errorf("unable to list directory in %v dir: %v", basedir, err)
	}

	imgCams := make([]string, 0)
	records := make([]string, 0)

	for _, dirItem := range dirItems {
		log.Debugf("process %v directory", dirItem)
		imgDir := path.Join(basedir, dirItem.Name(), camSubDir)
		imgs, err := ioutil.ReadDir(imgDir)
		if err != nil {
			return fmt.Errorf("unable to list cam images in directory %v: %v", imgDir, err)
		}

		for _, img := range imgs {
			idx, err := indexFromFile(img.Name())
			if err != nil {
				return fmt.Errorf("unable to find index in cam image name %v: %v", img.Name(), err)
			}
			log.Debugf("found image with index %v", idx)
			records = append(records, path.Join(basedir, dirItem.Name(), fmt.Sprintf(record.RecorNameFormat, idx)))
			imgCams = append(imgCams, path.Join(basedir, dirItem.Name(), camSubDir, img.Name()))
		}

	}

	if sliceSize > 0{
		imgCams, records, err = applySlice(imgCams, records, sliceSize)
	}

	content, err := buildArchiveContent(&imgCams, &records)
	if err != nil {
		return fmt.Errorf("unable to build archive: %v", err)
	}
	err = ioutil.WriteFile(archiveName, *content, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("unable to write archive content to disk: %v", err)
	}
	return nil
}

func applySlice(imgCams []string, records []string, sliceSize int) ([]string, []string, error) {
	// Add sliceSize images shift
	i := imgCams[:len(imgCams)-sliceSize]
	r := records[sliceSize:]

	return i, r, nil
}

var indexRegexp *regexp.Regexp

func init() {
	re, err := regexp.Compile("image_array_(?P<idx>[0-9]+)\\.jpg$")
	if err != nil {
		log.Fatalf("unable to compile regex: %v", err)
	}
	indexRegexp = re
}

func indexFromFile(fileName string) (string, error) {
	matches := findNamedMatches(indexRegexp, fileName)
	if matches["idx"] == "" {
		return "", fmt.Errorf("no index in filename")
	}
	return matches["idx"], nil
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

func buildArchiveContent(imgFiles *[]string, recordFiles *[]string) (*[]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	err := addJsonFiles(recordFiles, imgFiles, w)
	if err != nil {
		return nil, fmt.Errorf("unable to write json files in zip archive: %v", err)
	}

	err = addCamImages(imgFiles, w)
	if err != nil {
		return nil, fmt.Errorf("unable to cam files in zip archive: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to build archive: %v", err)
	}

	content, err := ioutil.ReadAll(buf)
	return &content, err
}

func addCamImages(imgFiles *[]string, w *zip.Writer) error {
	for _, img := range *imgFiles {
		imgContent, err := ioutil.ReadFile(img)
		if err != nil {
			return fmt.Errorf("unable to read img: %v", err)
		}

		_, imgName := path.Split(img)
		err = addToArchive(w, imgName, &imgContent)
		if err != nil {
			return fmt.Errorf("unable to create new img entry in archive: %v", err)
		}
	}
	return nil
}

func addJsonFiles(recordFiles *[]string, imgCam *[]string, w *zip.Writer) error {
	for idx, r := range *recordFiles {
		content, err := ioutil.ReadFile(r)
		if err != nil {
			return fmt.Errorf("unable to read json content: %v", err)
		}
		var rcd record.Record
		err = json.Unmarshal(content, &rcd)
		if err != nil {
			return fmt.Errorf("unable to unmarshal record: %v", err)
		}
		_, camName := path.Split((*imgCam)[idx])
		rcd.CamImageArray = camName

		recordBytes, err := json.Marshal(&rcd)
		if err != nil {
			return fmt.Errorf("unable to marshal %v record: %v", rcd, err)
		}

		_, recordName := path.Split(r)
		err = addToArchive(w, recordName, &recordBytes)
		if err != nil {
			return fmt.Errorf("unable to create new record in archive: %v", err)
		}
	}
	return nil
}

func addToArchive(w *zip.Writer, name string, content *[]byte) error {
	recordWriter, err := w.Create(name)
	if err != nil {
		return fmt.Errorf("unable to create new entry %v in archive: %v", name, err)
	}

	_, err = recordWriter.Write(*content)
	if err != nil {
		return fmt.Errorf("unable to add content in %v zip archive: %v", name, err)
	}
	return nil
}
