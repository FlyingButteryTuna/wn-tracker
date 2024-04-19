package json

import (
	"encoding/json"
	"os"

	"github.com/FlyingButterTuna/wn-tracker/novel"
)

func SaveNovelDataToFile(novelData novel.NovelData, filename string) error {
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
