package store

import (
	"go-web-exercises/internal/domain"
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
)

/*Interface for modify .json file of products
Initialization functions (verify that it is possible to read the file and modify it), 
Search (search for a specific product by id), 
Modify (update fields of a product by id) and 
Delete (delete a product by id).*/

func ReadJson() []domain.Product {
	var products []domain.Product
	jsonFile, err := os.Open(os.Getenv("PATH_JSON"))
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &products)
	return products
}

func WriteSlice(slice []domain.Product) error {
	jsonFile, err := os.OpenFile(os.Getenv("PATH_JSON"), os.O_RDWR | os.O_TRUNC, 0644)

	if err != nil {
		return errors.New("It was unable to open the file")
	}

	sliceToBytes, err := json.Marshal(slice)

	if err != nil {
		return errors.New("It was unable to write")
	}

	defer jsonFile.Close()
	jsonFile.Write(sliceToBytes)
	return nil
}