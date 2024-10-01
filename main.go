package main

import (
	"flag"
	"fmt"
	"kumquat/repository"
	"kumquat/store"
	"kumquat/template"
	"log/slog"
	"os"
)

func main() {
	var repo repository.Repository
	repo, err := repository.NewSQLiteRepository()
	if err != nil {
		slog.Error("Unable to create repository", "err", err)
		panic(err)
	}
	defer repo.Close() //nolint:errcheck

	inDir := flag.String("in", "sampledata", "directory path to read Kubernetes resources")
	flag.Parse()

	err = repository.LoadYAMLFromDirectoryTree(os.DirFS("."), *inDir, repo)
	if err != nil {
		slog.Error("Unable to load directory tree", "err", err)
		panic(err)
	}

	tplrs, err := repo.Query(
		/* sql */ `SELECT template.data AS tpl FROM "` + template.TemplateResourceType + `" AS template`,
	)
	if err != nil {
		slog.Error("Unable to find template", "err", err)
		panic(err)
	}

	// Process every Template
	templates := make([]*template.Template, 0, len(tplrs.Results))

	for _, tplrs := range tplrs.Results {
		tplres := tplrs["tpl"]
		t, err := template.NewTemplate(*tplres)

		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		templates = append(templates, t)
		fmt.Printf("Loaded Template %s\n", t.Name())
	}

	for _, t := range templates {
		o, err := t.Evaluate(repo)

		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		generateOutput(o)
	}
}

func generateOutput(o *template.TemplateOutput) {
	// for loop over data
	for i := 0; i < o.Output.ResourceCount(); i++ {
		out, err := o.Output.ResultString(i)

		if err != nil {
			panic(err)
		}

		fileName := o.FileNames[i]

		// use WriteFile function in store package to write the output to a file
		err = store.WriteToFile(fileName, "", out)

		if err != nil {
			panic(err)
		}
	}
}
