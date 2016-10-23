package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/adriwankenobi/comic/service"
	"github.com/adriwankenobi/comic/web"
	"io/ioutil"
)

func main() {

	// TODO: Custom usages for each flag
	convert := flag.Bool("convert", false, "Convert XLSX file to JSON")
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

	if *convert {
		out, errFlag := validateConvertFlags(*f, *o)
		if errFlag == nil {
			fmt.Printf("Converting from '%s' to '%s\n", *f, out)
			err = convertXLS(*f,
				fmt.Sprintf("%s/%s", out, web.ComicsFile),
				fmt.Sprintf("%s/%s", out, web.PhasesFile),
				fmt.Sprintf("%s/%s", out, web.FissuesFile),
			)
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

	if !*convert && !*update && !*folders {
		errFlag = errors.New("One these flags is mandatory: [-convert, -update, -folders]")
	}

	if errFlag != nil {
		fmt.Println(errFlag.Error())
		flag.PrintDefaults()
	}

	if err != nil {
		fmt.Println(err.Error())
	}

}

func validateConvertFlags(f, o string) (string, error) {
	if f == "" || o == "" {
		return "", errors.New("Input file and output path cannot be empty")
	}
	out := o
	if string(o[len(o)-1]) == "/" {
		out = o[:len(o)-1]
	}
	return out, nil
}

func convertXLS(f, comicsOut, phasesOut, issuesPhasesOut string) error {

	// Read XLS file
	comics, phases, issuesPhases, err := service.NewComicListFromXLSX(f)
	if err != nil {
		return err
	}

	// Write JSON files
	json, err := comics.ToJson()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(comicsOut, json, 0644)
	if err != nil {
		return err
	}
	json, err = phases.ToJson()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(phasesOut, json, 0644)
	if err != nil {
		return err
	}
	json, err = issuesPhases.ToJson()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(issuesPhasesOut, json, 0644)
	if err != nil {
		return err
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
