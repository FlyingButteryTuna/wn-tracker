package json

import (
	"encoding/json"
	"github.com/FlyingButterTuna/wn-tracker/parsers"
	"os"
)

func SaveNovelDataToFile(novelData parsers.NovelData, filename string) error {
	jsonData, err := json.MarshalIndent(novelData, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
