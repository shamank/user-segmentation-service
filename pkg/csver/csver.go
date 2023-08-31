package csver

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	assetsPath = "./assets/user_logs/"
)

type CSVer interface {
	CreateFile(date [][]interface{}) (string, error)
}

type CSVManager struct {
	assetsPath string
}

func NewCSVManager(assetsPath string) *CSVManager {
	return &CSVManager{
		assetsPath: assetsPath,
	}
}

func (cm *CSVManager) CreateFile(data [][]interface{}) (string, error) {
	name, err := generateRandomName()
	if err != nil {
		return "", err
	}

	filename := assetsPath + name + ".csv"

	err = os.MkdirAll(assetsPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	file, err := os.Create(filename)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Comma = ';'

	for _, row := range data {
		stringsRow := make([]string, len(row))

		for i, el := range row {
			switch t := el.(type) {
			case string:
				stringsRow[i] = t
			case int:
				stringsRow[i] = strconv.Itoa(t)
			case time.Time:
				stringsRow[i] = t.String()
			default:

				return "", fmt.Errorf("unknown types: %T", t)
			}

		}

		if err := writer.Write(stringsRow); err != nil {

			return "", err
		}
	}

	return name + ".csv", nil

}

func generateRandomName() (string, error) {
	// Создание буфера для случайных чисел
	buffer := make([]byte, 32)

	// Чтение случайных чисел в буфер
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
	}

	// Преобразование случайных чисел в строку шестнадцатеричного формата
	randomName := hex.EncodeToString(buffer)

	return randomName, nil
}
