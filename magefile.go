//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Builds the kumquat binary.
func Build() error {
	return sh.Run("go", "build", "-buildvcs=false")
}

// Runs all the examples in the examples directory, and compares them with the expected output.
func Examples() error {
	mg.Deps(Build)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	kcmd := path.Join(wd, "kumquat")

	examples, err := os.ReadDir("examples")
	if err != nil {
		return err
	}

	for _, example := range examples {
		if !example.IsDir() {
			continue
		}

		fmt.Printf("========== %s ==========\n", example.Name())

		cmd := exec.Command(kcmd, "-in", "input")
		exampleDir := path.Join("examples", example.Name())
		cmd.Dir = exampleDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error running example %s: %v\n", example.Name(), err)
			continue
		}

		// Compare the example output files with the files in the expected output directory
		err = sh.RunV("diff", "-r", path.Join(exampleDir, "output"), path.Join(exampleDir, "expected"))
		if err != nil {
			fmt.Printf("Error comparing example %s: %v\n", example.Name(), err)
		}
	}

	return nil
}

// Cleans all build, generated, and example outputs.
func Clean() error {
	mg.Deps(CleanExamples)

	err := os.RemoveAll("generated")
	if err != nil {
		return err
	}

	err = sh.Rm("kumquat")
	if err != nil {
		return err
	}

	return nil
}

// Cleans the example outputs.
func CleanExamples() error {
	examples, err := os.ReadDir("examples")
	if err != nil {
		return err
	}

	errs := make([]error, 0, len(examples))
	for _, example := range examples {
		if !example.IsDir() {
			continue
		}

		outputDir := path.Join("examples", example.Name(), "output")

		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			continue
		}

		fmt.Println(example.Name())
		err = os.RemoveAll(outputDir)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
