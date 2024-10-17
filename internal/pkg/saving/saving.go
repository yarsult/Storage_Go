package saving

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const rights = 0755

func WriteAtomic(path string, b []byte) error {
	logger, _ := zap.NewProduction()
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	tmpPathName := filepath.Join(dir, filename+".tmp")
	err := os.WriteFile(tmpPathName, b, rights)
	if err != nil {
		logger.Error("Failed to write JSON to file", zap.Error(err))
		return err
	}
	defer func() {
		os.Remove(tmpPathName)
	}()
	return os.Rename(tmpPathName, path)
}
