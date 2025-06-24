package files

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tfk70/hyprcircade/internal/logging"
)

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct the target path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Copy file
		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(srcFile, dstFile string, perm os.FileMode) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func TestReplaceInFile(t *testing.T) {
	logging.SetupLogger()

	testDataSrc := "./testdata/before-replace"
	testDataDst := "./testdata/_temp"
	testDataAfterReplace := "./testdata/after-replace"
	err := copyDir(testDataSrc, testDataDst)
	if err != nil {
		t.Fatalf("Error copying test data: %s", err.Error())
	}
	defer os.RemoveAll(testDataDst)

	anchor := "Anchor"

	anchorFileAfterReplace := testDataAfterReplace + "/anchor.yml"
	noAnchorFileAfterReplace := testDataAfterReplace + "/no-anchor.yml"
	oneLineNoAnchorFileAfterReplace := testDataAfterReplace + "/one-line-no-anchor"

	anchorFile := testDataDst + "/anchor.yml"
	noAnchorFile := testDataDst + "/no-anchor.yml"
	oneLineNoAnchorFile := testDataDst + "/one-line-no-anchor"

	err = ReplaceInFile(anchorFile, "edited", "EDITED", anchor)
	if err != nil {
		t.Fatalf("Error replacing in file %s: %s", anchor, err.Error())
	}

	anchorFileContent, err := os.ReadFile(anchorFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", anchorFile, err.Error())
	}
	anchorFileAfterReplaceContent, err := os.ReadFile(anchorFileAfterReplace)
	if err != nil {
		t.Fatalf("failed to read %s: %s", anchorFileAfterReplace, err.Error())
	}
	if diff := cmp.Diff(string(anchorFileContent), string(anchorFileAfterReplaceContent)); diff != "" {
		t.Errorf("Files %s and %s differ (-file1 +file2):\n%s", anchorFile, anchorFileAfterReplace, diff)
	}

	err = ReplaceInFile(noAnchorFile, "edited", "EDITED", "")
	if err != nil {
		t.Fatalf("Error replacing in file %s: %s", anchor, err.Error())
	}

	noAnchorFileContent, err := os.ReadFile(noAnchorFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", noAnchorFile, err.Error())
	}
	noAnchorFileAfterReplaceContent, err := os.ReadFile(noAnchorFileAfterReplace)
	if err != nil {
		t.Fatalf("failed to read %s: %s", noAnchorFileAfterReplace, err.Error())
	}
	if diff := cmp.Diff(string(noAnchorFileContent), string(noAnchorFileAfterReplaceContent)); diff != "" {
		t.Errorf("Files %s and %s differ (-file1 +file2):\n%s", noAnchorFile, noAnchorFileAfterReplace, diff)
	}

	err = ReplaceInFile(oneLineNoAnchorFile, "edit-me", "edited", "")
	if err != nil {
		t.Fatalf("Error replacing in file %s: %s", anchor, err.Error())
	}

	oneLineNoAnchorFileContent, err := os.ReadFile(oneLineNoAnchorFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %s", oneLineNoAnchorFile, err.Error())
	}
	oneLineNoAnchorFileAfterReplaceContent, err := os.ReadFile(oneLineNoAnchorFileAfterReplace)
	if err != nil {
		t.Fatalf("failed to read %s: %s", oneLineNoAnchorFileAfterReplace, err.Error())
	}
	if diff := cmp.Diff(string(oneLineNoAnchorFileContent), string(oneLineNoAnchorFileAfterReplaceContent)); diff != "" {
		t.Errorf("Files %s and %s differ (-file1 +file2):\n%s", oneLineNoAnchorFile, oneLineNoAnchorFileAfterReplace, diff)
	}
}
