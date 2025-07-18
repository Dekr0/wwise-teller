package wwise

import "testing"

func TestGenTranslation(t *testing.T) {
	if err := ReadTranslationCSV("./translation_table.csv", "./enum_table.csv"); err != nil {
		t.Fatal(err)
	}
}
