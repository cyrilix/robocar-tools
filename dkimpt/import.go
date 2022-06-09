package dkimpt

import (
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-tools/record"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

/* donkey import*/

var (
	camSubDir         = "cam"
	camIndexRegexp    *regexp.Regexp
	recordIndexRegexp *regexp.Regexp
)

func ImportDonkeyRecords(basedir string, destDir string) error {
	dirItems, err := ioutil.ReadDir(basedir)
	if err != nil {
		return fmt.Errorf("unable to list directory in %v dir: %v", basedir, err)
	}

	imgCams := make([]string, 0)
	records := make([]string, 0)

	for _, dirItem := range dirItems {
		zap.S().Debugf("process %v directory", dirItem)
		camDir := path.Join(destDir, dirItem.Name(), camSubDir)

		err := os.MkdirAll(camDir, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("unable to make dest directories %v: %v", camDir, err)
		}

		imgDir := path.Join(basedir, dirItem.Name(), camSubDir)
		imgs, err := ioutil.ReadDir(imgDir)
		if err != nil {
			return fmt.Errorf("unable to list cam images in directory %v: %v", imgDir, err)
		}

		for _, img := range imgs {
			idx, err := indexFromFile(camIndexRegexp, img.Name())
			if err != nil {
				return fmt.Errorf("unable to find index in cam image name %v: %v", img.Name(), err)
			}
			zap.S().Debugf("found image with index %v", idx)
			records = append(records, path.Join(basedir, dirItem.Name(), fmt.Sprintf(record.FileNameFormat, idx)))
			imgCams = append(imgCams, path.Join(basedir, dirItem.Name(), camSubDir, img.Name()))
		}

		err = copyToDestdir(destDir, dirItem.Name(), &imgCams, &records)
		if err != nil {
			zap.S().Warnf("unable to copy files from %v to %v: %v", path.Join(basedir, dirItem.Name()), destDir, err)
			continue
		}
	}

	return nil
}

func init() {
	re, err := regexp.Compile("image_array_(?P<idx>[0-9]+)\\.jpg$")
	if err != nil {
		log.Fatalf("unable to compile regex: %v", err)
	}
	camIndexRegexp = re

	re, err = regexp.Compile("record_(?P<idx>[0-9]+)\\.json$")
	if err != nil {
		log.Fatalf("unable to compile regex: %v", err)
	}
	recordIndexRegexp = re
}

func indexFromFile(regex *regexp.Regexp, fileName string) (string, error) {
	matches := findNamedMatches(regex, fileName)
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

func copyToDestdir(destdir, dirItem string, imgFiles *[]string, recordFiles *[]string) error {
	err := copyJsonFiles(destdir, dirItem, recordFiles)
	if err != nil {
		return fmt.Errorf("unable to copy json files in %v directory: %v", destdir, err)
	}

	err = copyCamImages(destdir, dirItem, imgFiles)
	if err != nil {
		return fmt.Errorf("unable to copy cam files in %v director: %v", destdir, err)
	}

	return nil
}

func copyCamImages(destdir, dirItem string, imgFiles *[]string) error {
	for _, img := range *imgFiles {
		imgContent, err := ioutil.ReadFile(img)
		if err != nil {
			return fmt.Errorf("unable to read img: %v", err)
		}

		idx, err := indexFromFile(camIndexRegexp, img)
		if err != nil {
			zap.S().Warnf("unable to extract idx from filename %v: %v", img, err)
			continue
		}
		imgName := path.Join(destdir, dirItem, camSubDir, fmt.Sprintf("image_array_%s_%s.jpg", dirItem, idx))

		err = ioutil.WriteFile(imgName, imgContent, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("unable to write image %v: %v", imgName, err)
		}
	}
	return nil
}

func copyJsonFiles(destdir, dirItem string, recordFiles *[]string) error {
	for _, r := range *recordFiles {
		content, err := ioutil.ReadFile(r)
		if err != nil {
			return fmt.Errorf("unable to read json content: %v", err)
		}
		idx, err := indexFromFile(recordIndexRegexp, r)
		if err != nil {
			zap.S().Warnf("unable to extract idx from filename %v: %v", r, err)
			continue
		}

		var rcd record.Record
		err = json.Unmarshal(content, &rcd)
		if err != nil {
			return fmt.Errorf("unable to unmarshal record: %v", err)
		}
		camName := fmt.Sprintf("image_array_%s_%s.jpg", dirItem, idx)
		rcd.CamImageArray = path.Join(camSubDir, camName)

		recordBytes, err := json.Marshal(&rcd)
		if err != nil {
			return fmt.Errorf("unable to marshal %v record: %v", rcd, err)
		}

		recordFileName := path.Join(destdir, dirItem, fmt.Sprintf("record_%s_%s.json", dirItem, idx))
		err = ioutil.WriteFile(recordFileName, recordBytes, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("unable to write json record %v: %v", recordFileName, err)
		}
	}
	return nil
}
