package main

// When adding new package to template, remove mod file and add name to the slice
// Packages must be specified in the order of their dependencies
// Replace organization name with {{.orgName}}, package repository name with {{.pkgRepoName}}

var pkgOrder = []string{
	"errors",
	"log",
	"env",
	"app",
	"models",
	"mongo",
	"nats",
	"redis",
	"rest",
	"utils",
}
