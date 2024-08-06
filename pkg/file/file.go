package file

import (
	"errors"
	"io"
	"net/http"
)

func GetFile(r *http.Request, key string) (fileName string, data []byte, err error) {
	file, header, err := r.FormFile("avatar")
	if err != nil {
		return "", nil, err
	}
	defer file.Close()
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", nil, errors.New("ошибка при чтении файла аватара")
	}
	//if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
	//	response.Fail(w, http.StatusBadRequest, "Недопустимый тип файла для аватара")
	//	return
	//}
	if header.Size > 5*1024*1024 {
		return "", nil, errors.New("размер файла превышен")
	}

	avatarData, err := io.ReadAll(file)
	if err != nil {
		return "", nil, errors.New("ошибка при чтении файла аватара")
	}
	return header.Filename, avatarData, nil
}
