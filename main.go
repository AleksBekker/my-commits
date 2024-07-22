package main

import (
	"flag"
	"fmt"
	"github.com/AleksBekker/my-commits/checker"
	"github.com/AleksBekker/my-commits/colors"
	"github.com/AleksBekker/my-commits/git"
	"os"
	"regexp"
)

// pattern to validate emails
const emailPattern = `(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

const (
	error_ = iota
	warn
	info
	all
)

var (
	infoColor    colors.Colors = []string{colors.CyanFg}
	warningColor colors.Colors = []string{colors.YellowFg}
	errorColor   colors.Colors = []string{colors.RedFg}
)

// Arguments represents all possible command line arguments
type Arguments struct {
	dir          string
	email        string
	since        string
	targetRemote string
	checkLinks   bool
	verbosity    int
}

// clArgs represents parsed command line arguments
var clArgs = parseArguments()

func main() {

	if isGit, _ := git.IsGitRepo(clArgs.dir); !isGit {
		_, _ = errorPrintf("'%s' is not part of a Git directory\n", clArgs.dir)
		os.Exit(1)
	}

	params, err := git.GetLogParams(clArgs.dir, clArgs.targetRemote, clArgs.since)
	if err != nil {
		_, _ = errorPrintf("something went wrong while getting log parameters: %s\n", err.Error())
		os.Exit(1)
	}

	_, _ = infoPrintf("Parameters:\n")
	_, _ = infoPrintf("    Dir: %s\n", params.Dir)
	_, _ = infoPrintf("    Since: %s\n", params.Since)
	_, _ = infoPrintf("    Email: %s\n", params.Email)
	_, _ = infoPrintf("    Commit Link: %s\n", params.CommitLink)

	if clArgs.email != "" {
		params.Email = clArgs.email
	}

	if !regexp.MustCompile(emailPattern).Match([]byte(params.Email)) {
		_, _ = warnPrintf("email '%s' doesn't appear to be a valid email\n", params.Email)
	}

	links, _ := git.GetCommits(params)

	if clArgs.checkLinks {
		results := checker.CheckLinks(links)
		for _, result := range results {
			switch result.Status {
			case 200:
				fmt.Println(result.Link)
			case 404:
				_, _ = infoPrintf("'%s' not found\n", result.Link)
			case 429:
				_, _ = warnPrintf("'%s' cannot be checked due to rate limiting -- please try again later\n", result.Link)
			default:
				_, _ = warnPrintf("link '%s' returns unexpected status %d\n", result.Link, result.Status)
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
		args.verbosity = error_
	} else if verbose {
		args.verbosity = all
	} else {
		args.verbosity = warn
	}

	return args
}

func infoPrintf(format string, args ...any) (int, error) {
	if clArgs.verbosity < info {
		return 0, nil
	}

	return infoColor.Fprintf(os.Stderr, "INFO: %s", fmt.Sprintf(format, args...))
}

func warnPrintf(format string, args ...any) (int, error) {
	if clArgs.verbosity < warn {
		return 0, nil
	}

	return warningColor.Fprintf(os.Stderr, "WARNING: %s", fmt.Sprintf(format, args...))
}

func errorPrintf(format string, args ...any) (int, error) {
	if clArgs.verbosity < error_ {
		return 0, nil
	}

	return errorColor.Fprintf(os.Stderr, "ERROR: %s", fmt.Sprintf(format, args...))
}
