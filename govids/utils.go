package govids

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ValidatePath helper to validate if a path exists
func ValidatePath(p string) string {
	v, err := filepath.Abs(p)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(v); os.IsNotExist(err) {
		log.Fatal(err)
	}

	return v
}

// ReadJSON process json video file and return array of Vidoes
func ReadJSON(in string) []byte {
	fmt.Println("Reading JSON")

	jsonFile, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}
