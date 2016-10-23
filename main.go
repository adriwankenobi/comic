package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"io/ioutil"
)

func main() {

	// TODO: Custom usages for each flag
	generate := flag.Bool("generate", false, "Generate JSON files from XLSX file")
	update := flag.Bool("update", false, "Update XLSX file with some info from MARVEL API")
	folders := flag.Bool("folders", false, "Create folders structure")
	f := flag.String("f", "", "XSLX file to read")
	o := flag.String("o", "", "Path to output")
	start := flag.Int("start", -1, "Start year to find comics")
	end := flag.Int("end", -1, "End year to find comics")
	mPubKey := flag.String("mpubkey", "", "MARVEL API public key")
	mPriKey := flag.String("mprikey", "", "MARVEL API private key")
	flag.Parse()

	var err error
	var errFlag error

	if *generate {
		out, errFlag := validateGenerateFlags(*f, *o)
		if errFlag == nil {
			fmt.Printf("Generating from '%s' to '%s'\n", *f, out)
			err = generateJSON(*f, out)
		}
	}

	if *update {
		errFlag = validateUpdateFlags(*f, *start, *end, *mPubKey, *mPriKey)
		if errFlag == nil {
			fmt.Printf("Updating '%s'\n", *f)
			err = updateXLS(*f, *start, *end, *mPubKey, *mPriKey)
		}
	}

	if *folders {
		out, errFlag := validateFoldersFlags(*f, *o)
		if errFlag == nil {
			fmt.Printf("Creating folders from '%s' in '%s\n", *f, out)
			err = createFolders(*f, out)
		}
	}

	if !*generate && !*update && !*folders {
		errFlag = errors.New("One these flags is mandatory: [-generate, -update, -folders]")
	}

	if errFlag != nil {
		fmt.Println(errFlag.Error())
		flag.PrintDefaults()
	}

	if err != nil {
		fmt.Println(err.Error())
	}

}

func validateGenerateFlags(f, o string) (string, error) {
	if f == "" || o == "" {
		return "", errors.New("Input file and output path cannot be empty")
	}
	out := o
	if string(o[len(o)-1]) == "/" {
		out = o[:len(o)-1]
	}
	return out, nil
}

func generateJSON(f, out string) error {
	// Read XLS file
	err := service.JsonGenerator(f, out)
	if err != nil {
		return err
	}

	// Write JSON files
	for key, value := range service.Datastore {
		json, err := value.ToJson()
		if err != nil {
			return err
		}
		output := fmt.Sprintf("%s/%s.json", out, key)
		err = ioutil.WriteFile(output, json, 0644)
		if err != nil {
			return err
		}
	}
	fmt.Println("Done!")
	return nil
}

func validateUpdateFlags(f string, start, end int, mPubKey, mPriKey string) error {
	if f == "" || start == -1 || end == -1 || mPubKey == "" || mPriKey == "" {
		return errors.New("Input file, start and end cannot be empty\nMARVEL public and private keys are needed")
	}
	return nil
}

func updateXLS(f string, start, end int, mPubKey, mPriKey string) error {

	// Update XLS file
	err := service.UpdateXLSX(f, start, end, mPubKey, mPriKey)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

func validateFoldersFlags(f, o string) (string, error) {
	if f == "" || o == "" {
		return "", errors.New("Input file and output path cannot be empty")
	}
	out := o
	if string(o[len(o)-1]) == "/" {
		out = o[:len(o)-1]
	}
	return out, nil
}

func createFolders(f, o string) error {

	// Create folders structure
	err := service.CreateFolders(f, o)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}
