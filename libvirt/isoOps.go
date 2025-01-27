package libvirt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"udp_iaas/types"
)

func ListISOs() ([]types.ISO, error) {
	if err := os.MkdirAll(isoStoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create iso directory: %w", err)
	}

	entries, err := os.ReadDir(isoStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read iso directory: %w", err)
	}

	var isos []types.ISO
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".iso") {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			isos = append(isos, types.ISO{
				Name: entry.Name(),
				Size: info.Size(),
			})
		}
	}

	return isos, nil
}

func SaveISO(fileName string, file io.Reader) error {
	if !strings.HasSuffix(strings.ToLower(fileName), ".iso") {
		return fmt.Errorf("file must have .iso extension")
	}

	if err := os.MkdirAll(isoStoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create iso directory: %w", err)
	}

	dst, err := os.Create(filepath.Join(isoStoragePath, fileName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}