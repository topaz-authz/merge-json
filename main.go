package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"dario.cat/mergo"
	"github.com/samber/lo"
)

func main() {
	fmt.Println("merge-json")

	output := flag.String("output", "", "output filename")

	flag.Parse()

	fmt.Println("inputs:")
	lo.ForEach(flag.Args(), func(x string, _ int) {
		fmt.Printf("=> %s\n", x)
	})

	fmt.Println("output:", *output)

	if err := merge(flag.Args(), *output); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func merge(files []string, outfile string) error {
	var (
		src map[string]any
		dst = make(map[string]any)
		err error
	)

	for _, file := range files {
		fmt.Printf("<= %s\n", file)

		src, err = loadFile(file)
		if err != nil {
			return err
		}

		if err := mergo.Merge(&dst, src); err != nil {
			return err
		}
	}

	w, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(true)

	fmt.Printf("=> %s\n", outfile)

	err = enc.Encode(dst)

	return err
}

func loadFile(filePath string) (map[string]any, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data map[string]any

	dec := json.NewDecoder(f)

	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
