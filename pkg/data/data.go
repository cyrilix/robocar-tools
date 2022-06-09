package data

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-tools/record"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var camSubDir = "cam"

func WriteArchive(basedir string, archiveName string, sliceSize int, imgWidth, imgHeight int, horizon int, flipImages bool) error {
	content, err := BuildArchive(basedir, sliceSize, imgWidth, imgHeight, horizon, flipImages)
	if err != nil {
		return fmt.Errorf("unable to build archive: %w", err)
	}
	if err != nil {
		return fmt.Errorf("unable to build archive: %w", err)
	}
	err = ioutil.WriteFile(archiveName, content, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("unable to write archive content to disk: %w", err)
	}
	return nil
}

func BuildArchive(basedir string, sliceSize int, imgWidth, imgHeight int, horizon int, flipImages bool) ([]byte, error) {
	l := zap.S()
	l.Infof("build zip archive from %s\n", basedir)
	dirItems, err := ioutil.ReadDir(basedir)
	if err != nil {
		return nil, fmt.Errorf("unable to list directory in %v dir: %w", basedir, err)
	}

	imgCams := make([]string, 0)
	records := make([]string, 0)

	for _, dirItem := range dirItems {
		l.Infof("process %v directory", dirItem.Name())
		imgDir := path.Join(basedir, dirItem.Name(), camSubDir)
		imgs, err := ioutil.ReadDir(imgDir)
		if err != nil {
			return nil, fmt.Errorf("unable to list cam images in directory %v: %w", imgDir, err)
		}

		for _, img := range imgs {
			idx, err := indexFromFile(img.Name())
			if err != nil {
				return nil, fmt.Errorf("unable to find index in cam image name %v: %w", img.Name(), err)
			}
			l.Debugf("found image with index %v", idx)
			records = append(records, path.Join(basedir, dirItem.Name(), fmt.Sprintf(record.FileNameFormat, idx)))
			imgCams = append(imgCams, path.Join(basedir, dirItem.Name(), camSubDir, img.Name()))
		}
	}

	if sliceSize > 0 {
		imgCams, records, err = applySlice(imgCams, records, sliceSize)
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)

	err = buildArchiveContent(w, imgCams, records, imgWidth, imgHeight, horizon, false)
	if err != nil {
		return nil, fmt.Errorf("unable to build archive: %w", err)
	}
	if flipImages {
		err = buildArchiveContent(w, imgCams, records, imgWidth, imgHeight, horizon, true)
		if err != nil {
			return nil, fmt.Errorf("unable to build archive: %w", err)
		}
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to close zip archive: %w", err)
	}
	content, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("unable to generate archive content: %w", err)
	}
	l.Info("archive built\n")
	return content, nil
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
		zap.S().Fatalf("unable to compile regex: %v", err)
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

func buildArchiveContent(w *zip.Writer, imgFiles []string, recordFiles []string, imgWidth, imgHeight int, horizon int, withFlipImages bool) error {
	err := addJsonFiles(recordFiles, imgFiles, withFlipImages, w)
	if err != nil {
		return fmt.Errorf("unable to write json files in zip archive: %w", err)
	}

	err = addCamImages(imgFiles, withFlipImages, w, imgWidth, imgHeight, horizon)
	if err != nil {
		return fmt.Errorf("unable to cam files in zip archive: %w", err)
	}

	return err
}

func addCamImages(imgFiles []string, flipImage bool, w *zip.Writer, imgWidth, imgHeight int, horizon int) error {
	for _, im := range imgFiles {
		imgContent, err := ioutil.ReadFile(im)
		if err != nil {
			return fmt.Errorf("unable to read img: %w", err)
		}

		_, imgName := path.Split(im)
		if flipImage || imgWidth > 0 && imgHeight > 0 || horizon > 0 {
			img, _, err := image.Decode(bytes.NewReader(imgContent))
			if err != nil {
				zap.S().Fatalf("unable to decode jpeg image: %v", err)
			}

			if imgWidth > 0 && imgHeight > 0 {
				bounds := img.Bounds()
				if bounds.Dx() != imgWidth || bounds.Dy() != imgWidth {
					zap.S().Debugf("resize image %v from %dx%d to %dx%d", im, bounds.Dx(), bounds.Dy(), imgWidth, imgHeight)
					img = imaging.Resize(img, imgWidth, imgHeight, imaging.NearestNeighbor)
				}
			}
			if flipImage {
				img = imaging.FlipH(img)
				imgName = fmt.Sprintf("flip_%s", imgName)
			}
			if horizon > 0 {
				img = imaging.Crop(img, image.Rect(0, horizon, img.Bounds().Dx(), img.Bounds().Dy()))
			}
			var bytesBuff bytes.Buffer
			err = jpeg.Encode(&bytesBuff, img, nil)
			imgContent = bytesBuff.Bytes()
		}

		err = addToArchive(w, imgName, imgContent)
		if err != nil {
			return fmt.Errorf("unable to create new img entry in archive: %w", err)
		}
	}
	return nil
}

func addJsonFiles(recordFiles []string, imgCam []string, flipImage bool, w *zip.Writer) error {
	for idx, r := range recordFiles {
		content, err := ioutil.ReadFile(r)
		if err != nil {
			return fmt.Errorf("unable to read json content: %w", err)
		}
		var rcd record.Record
		err = json.Unmarshal(content, &rcd)
		if err != nil {
			return fmt.Errorf("unable to unmarshal record: %w", err)
		}
		_, camName := path.Split((imgCam)[idx])

		if flipImage {
			rcd.UserAngle = rcd.UserAngle * -1
			rcd.CamImageArray = fmt.Sprintf("flip_%s", camName)
		}else {
			rcd.CamImageArray = camName
		}

		recordBytes, err := json.Marshal(&rcd)
		if err != nil {
			return fmt.Errorf("unable to marshal %v record: %w", rcd, err)
		}

		_, recordName := path.Split(r)
		if flipImage {
			recordName = strings.ReplaceAll(recordName, "record", "record_flip")
		}
		err = addToArchive(w, recordName, recordBytes)
		if err != nil {
			return fmt.Errorf("unable to create new record in archive: %w", err)
		}
	}
	return nil
}

func addToArchive(w *zip.Writer, name string, content []byte) error {
	recordWriter, err := w.Create(name)
	if err != nil {
		return fmt.Errorf("unable to create new entry %v in archive: %w", name, err)
	}

	_, err = recordWriter.Write(content)
	if err != nil {
		return fmt.Errorf("unable to add content in %v zip archive: %w", name, err)
	}
	return nil
}
