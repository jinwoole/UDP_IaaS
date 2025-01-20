package system

import (
	"os"
	"path/filepath"
)

func (d *Data) GetISOFiles() ([]ISOInfo, error) {
	var filesInfo []ISOInfo
	err := filepath.Walk("iso", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesInfo = append(filesInfo, ISOInfo{
				Name: info.Name(),
				Path: path,
			})
		}
		return nil
	})
	return filesInfo, err
}
