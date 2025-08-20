package file_utils

import (
	"io"
	"net/http"
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

// SaveUploadedFile получает файл из multipart-формы по ключу и сохраняет его по указанному пути.
// Возвращает true, если файл был найден и сохранен, и false, если файл не был найден.
func SaveUploadedFile(r *http.Request, formKey string, destinationPath string) (bool, error) {
	file, _, err := r.FormFile(formKey)
	if err != nil {
		// Если файл не найден, это не ошибка, просто пропускаем.
		if err == http.ErrMissingFile {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	if err := SaveFileFromReader(file, destinationPath); err != nil {
		return false, err
	}

	return true, nil
}
