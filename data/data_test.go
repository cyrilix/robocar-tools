package data

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/cyrilix/robocar-tools/record"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestBuildArchive(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "buildarchive")
	if err != nil {
		t.Fatalf("unable to make tmpdir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			log.Warnf("unable to remove tempdir %v: %v", tmpDir, err)
		}
	}()

	archive := path.Join(tmpDir, "train.zip")

	expectedRecordFiles, expectedImgFiles := expectedFiles()

	err = BuildArchive("testdata", archive, 0)
	if err != nil {
		t.Errorf("unable to build archive: %v", err)
	}

	r, err := zip.OpenReader(archive)
	if err != nil {
		t.Errorf("unable to read archive, %v", err)
	}
	defer r.Close()

	if len(r.File) != len(expectedImgFiles)+len(expectedRecordFiles) {
		t.Errorf("bad number of files in archive: %v, wants %v", len(r.File), len(expectedImgFiles)+len(expectedRecordFiles))
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		filename := f.Name
		if filename[len(filename)-4:] == "json" {
			expectedRecordFiles[filename] = true
			expectedtImgName := strings.Replace(filename, "record", "cam-image_array", 1)
			expectedtImgName = strings.Replace(expectedtImgName, "json", "jpg", 1)
			checkJsonContent(t, f, expectedtImgName)
			continue
		}
		if filename[len(filename)-3:] == "jpg" {
			expectedImgFiles[filename] = true
			continue
		}
		t.Errorf("unexpected file in archive: %v", filename)
	}

	checkAllFilesAreFoundInArchive(expectedRecordFiles, t, expectedImgFiles)
}

func checkAllFilesAreFoundInArchive(expectedRecordFiles map[string]bool, t *testing.T, expectedImgFiles map[string]bool) {
	for f, found := range expectedRecordFiles {
		if !found {
			t.Errorf("%v not found in archive", f)
		}
	}
	for f, found := range expectedImgFiles {
		if !found {
			t.Errorf("%v not found in archive", f)
		}
	}
}

func checkJsonContent(t *testing.T, f *zip.File, expectedCamImage string) {
	rc, err := f.Open()
	if err != nil {
		t.Errorf("unable to read file content of %v: %v", f.Name, err)
	}
	defer rc.Close()

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Errorf("%v has invalid json content: %v", f.Name, err)
	}
	var rcd record.Record
	err = json.Unmarshal(content, &rcd)
	if err != nil {
		t.Errorf("unable to unmarshal json content of%v: %v", f.Name, err)
	}

	if rcd.CamImageArray != expectedCamImage {
		t.Errorf("record %v: invalid image ref: %v, wants %v", f.Name, rcd.CamImageArray, expectedCamImage)
	}
	if rcd.UserAngle == 0. {
		t.Errorf("record %v: user angle has not been initialised", f.Name)
	}
}

func expectedFiles() (map[string]bool, map[string]bool) {
	expectedRecordFiles := make(map[string]bool)
	expectedImgFiles := make(map[string]bool)
	for i := 1; i <= 8; i++ {
		expectedRecordFiles[fmt.Sprintf("record_%07d.json", i)] = false
		expectedImgFiles[fmt.Sprintf("cam-image_array_%07d.jpg", i)] = false
	}
	for i := 101; i <= 106; i++ {
		expectedRecordFiles[fmt.Sprintf("record_%07d.json", i)] = false
		expectedImgFiles[fmt.Sprintf("cam-image_array_%07d.jpg", i)] = false
	}
	return expectedRecordFiles, expectedImgFiles
}
