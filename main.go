package main

import (
	"fmt"
	"golang/prepare"
	"golang/utils"
	"os"
)

func main() {
	var prjName, orgName, pkgRepoName string

	fmt.Println("Creating new project barebone...")
	fmt.Println()
	fmt.Println("****** THINGS MUST BE DONE ******")
	fmt.Println("1. Create an organization on GitHub")
	fmt.Println("2. Create a packages repository on GitHub")
	fmt.Println()

	fmt.Println("What is the name of new project?")
	if _, err := fmt.Scanln(&prjName); err != nil {
		panic(err)
	}

	fmt.Println("What is the name of the organization on GitHub?")
	if _, err := fmt.Scanln(&orgName); err != nil {
		panic(err)
	}

	fmt.Println("What is the name of the packages repository on GitHub?")
	if _, err := fmt.Scanln(&pkgRepoName); err != nil {
		panic(err)
	}

	outDir := "./output"

	if err := os.RemoveAll(outDir); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(outDir, 0777); err != nil {
		panic(err)
	}

	if err := utils.CopyDir("./template", outDir); err != nil {
		panic(err)
	}

	if err := prepare.Pkg(outDir, prjName, orgName, pkgRepoName, pkgOrder); err != nil {
		panic(err)
	}

	if err := prepare.Generators(outDir, orgName); err != nil {
		panic(err)
	}

	if err := prepare.DevHelper(outDir, orgName); err != nil {
		panic(err)
	}
}
