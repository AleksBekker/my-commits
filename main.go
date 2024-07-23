package main

import (
	"flag"
	"fmt"
	"github.com/AleksBekker/my-commits/checker"
	"github.com/AleksBekker/my-commits/git"
	"github.com/AleksBekker/my-commits/logger"
	"os"
	"regexp"
)

// pattern to validate emails
const emailPattern = `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

// Arguments represents all possible command line arguments
type Arguments struct {
	dir          string
	email        string
	since        string
	targetRemote string
	checkLinks   bool
	verbosity    int
}

func main() {

	var clArgs = parseArguments()
	l := logger.Default().SetVerbosity(clArgs.verbosity)

	if isGit, _ := git.IsGitRepo(clArgs.dir); !isGit {
		_, _ = l.Fatalf("'%s' is not part of a Git directory\n", clArgs.dir)
		os.Exit(1)
	}

	params, err := git.GetLogParams(clArgs.dir, clArgs.targetRemote, clArgs.since)
	if err != nil {
		_, _ = l.Fatalf("something went wrong while getting log parameters: %s\n", err.Error())
		os.Exit(1)
	}

	// TODO: make a l.Debugf function
	//_, _ = l.Infof("Parameters: { Dir: '%s', Since: '%s', Email: '%s', Commit Link: '%s' }\n",
	//	params.Dir, params.Since, params.Email, params.CommitLink)

	if clArgs.email != "" {
		params.Email = clArgs.email
	}

	if !regexp.MustCompile(emailPattern).Match([]byte(params.Email)) {
		_, _ = l.Warnf("email '%s' doesn't appear to be a valid email\n", params.Email)
	}

	links, _ := git.GetCommits(params)

	if clArgs.checkLinks {
		results := checker.CheckLinks(links)
		for _, result := range results {
			switch result.Status {
			case 200:
				fmt.Println(result.Link)
			case 404:
				_, _ = l.Infof("'%s' not found\n", result.Link)
			case 429:
				_, _ = l.Warnf("'%s' cannot be checked due to rate limiting -- please try again later\n", result.Link)
			default:
				_, _ = l.Warnf("link '%s' returns unexpected status %d\n", result.Link, result.Status)
			}
		}
	} else {
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

func parseArguments() Arguments {
	args := Arguments{}
	flag.StringVar(&args.dir, "dir", ".", "Directory of the git repository")
	flag.StringVar(&args.email, "email", "", "The email associated with commits")
	flag.StringVar(&args.since, "since", "1970-01-01", "Furthest back the printed commits will go")
	flag.StringVar(&args.targetRemote, "remote", "origin", "The remote to use for link generation")
	flag.BoolVar(&args.checkLinks, "check", false, "Check if commit links are valid")

	var verbose, quiet bool
	flag.BoolVar(&verbose, "v", false, "Print additional information (overridden by -q)")
	flag.BoolVar(&quiet, "q", false, "Print no addition information (overrides -v)")

	flag.Parse()

	if quiet {
		args.verbosity = logger.Err
	} else if verbose {
		args.verbosity = logger.All
	} else {
		args.verbosity = logger.Warn
	}

	return args
}
