package file_utils

import (
	"io"
	"os"
)

// SaveFileFromReader читает данные из io.Reader и сохраняет их в файл по указанному пути.
func SaveFileFromReader(src io.Reader, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
