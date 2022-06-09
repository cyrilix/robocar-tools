package dkimpt

import (
	"encoding/json"
	record2 "github.com/cyrilix/robocar-tools/record"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestImportDonkeyRecords(t *testing.T) {
	destDir, err := ioutil.TempDir("", "test-import")
	if err != nil {
		t.Errorf("unable to generate destdir for test")
	}
	defer func() {
		err := os.RemoveAll(destDir)
		zap.S().Errorf("unable to delete tmpdir %v: %v", destDir, err)
	}()

	err = ImportDonkeyRecords("testdata", destDir)
	if err != nil {
		t.Errorf("unable to import files: %v", err)
	}

	cases := []struct {
		setDir          string
		expectedDir     int
		expectedRecords int
	}{
		{"20191012_111416", 1, 7},
		{"20191012_122633", 1, 8},
	}

	for _, c := range cases {

		files, err := ioutil.ReadDir(path.Join(destDir, c.setDir))
		if err != nil {
			t.Errorf("[%v] unable to list %v content: %v", c.setDir, destDir, err)
		}

		dirNames := make([]string, 0)
		fileNames := make([]string, 0)
		for _, f := range files {
			if f.IsDir() {
				dirNames = append(dirNames, f.Name())
			} else {
				fileNames = append(fileNames, f.Name())
			}
		}

		if len(dirNames) != c.expectedDir {
			t.Errorf("[%v] %v dirs found, wants %v", c.setDir, len(dirNames), 1)
		}
		if len(fileNames) != c.expectedRecords {
			t.Errorf("[%v] %v files found, wants %v", c.setDir, len(fileNames), 15)
		}

		for _, jsonFile := range fileNames {
			content, err := ioutil.ReadFile(path.Join(destDir, c.setDir, jsonFile))
			if err != nil {
				t.Errorf("[%v] unable to read file content of %v: %v", c.setDir, jsonFile, err)
				continue
			}
			var record record2.Record
			err = json.Unmarshal(content, &record)
			if err != nil {
				t.Errorf("[%v] unable to unmarshal file %v with content %v: %v", c.setDir, jsonFile, content, err)
				continue
			}
			if record.UserAngle == 0. {
				t.Errorf("[%v] user_angle field not initialized for %v file: %v", c.setDir, jsonFile, content)
			}
			camPath := path.Join(destDir, c.setDir, record.CamImageArray)
			camFile, err := os.Stat(camPath)
			if os.IsNotExist(err) {
				t.Errorf("[%v] cam image %v doesn't exist for record %v: %v", c.setDir, camPath, jsonFile, err)
				continue
			}
			if camFile.Size() == 0 {
				t.Errorf("[%v] cam image %v is empty", c.setDir, camFile.Name())
			}
		}
	}
}
