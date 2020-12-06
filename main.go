package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

const GitDir = ".git"

const GitConf = ".git/config"

func main() {
	app := cli.NewApp()
	app.Name = "gitcon"
	app.Usage = ".git/config management cli tool"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "ssh",
			Usage: "Switch to ssh repository",
			Action: func(c *cli.Context) error {
				newConfig := BuildConfig("ssh")
				fmt.Println(newConfig)
				err := ioutil.WriteFile(GitConf, []byte(newConfig), 0600)
				if err != nil {
					panic(err)
				}

				return nil
			},
		},
		{
			Name:  "https",
			Usage: "Switch to https repository",
			Action: func(c *cli.Context) error {
				newConfig := BuildConfig("https")
				fmt.Println(newConfig)
				err := ioutil.WriteFile(GitConf, []byte(newConfig), 0600)
				if err != nil {
					panic(err)
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func findGitCohnfig() bool {
	if f, err := os.Stat(GitDir); os.IsNotExist(err) || f.IsDir() {
		if conf, err := os.Stat(GitConf); os.IsNotExist(err) || !conf.IsDir() {
			return true
		}

		return false
	}

	return false
}

func ReplaceHTTPStoSSH(scanner *bufio.Scanner) string {
	newConfig := make([]string, 0)
	r := regexp.MustCompile(`url \= https://github\.com/([\w\d]+)/([\w\d]+)`)

	for scanner.Scan() {
		line := scanner.Bytes()
		lineStr := string(line)

		if r.FindString(lineStr) != "" {
			matches := r.FindSubmatch(line)
			old := string(matches[0])
			account := string(matches[1])
			repository := string(matches[2])
			sshURL := "url = git@github.com::" + account + "/" + repository + ".git"
			lineStr = strings.Replace(lineStr, old, sshURL, 1)
			newConfig = append(newConfig, lineStr)
		} else {
			newConfig = append(newConfig, lineStr)
		}
	}

	result := strings.Join(newConfig, "\n")

	return result
}

func ReplaceSSHtoHTTPS(scanner *bufio.Scanner) string {
	newConfig := make([]string, 0)
	r := regexp.MustCompile(`url \= git\@github\.com\:\:([\w\d]+)/([\w\d]+)\.git`)

	for scanner.Scan() {
		line := scanner.Bytes()
		lineStr := string(line)

		if r.FindString(lineStr) != "" {
			matches := r.FindSubmatch(line)
			old := string(matches[0])
			account := string(matches[1])
			repository := string(matches[2])
			// url = https://github.com/tMinamiii/woke
			httpsURL := "url = https://github.com/" + account + "/" + repository
			lineStr = strings.Replace(lineStr, old, httpsURL, 1)
			newConfig = append(newConfig, lineStr)
		} else {
			newConfig = append(newConfig, lineStr)
		}
	}

	result := strings.Join(newConfig, "\n")

	return result
}

func BuildConfig(mode string) string {
	newConfig := make([]string, 0)

	if !findGitCohnfig() {
		panic("err")
	}

	data, _ := os.Open(GitConf)
	defer data.Close()
	scanner := bufio.NewScanner(data)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "[remote \"origin\"]" {
			if mode == "ssh" {
				result := ReplaceHTTPStoSSH(scanner)

				newConfig = append(newConfig, line)
				newConfig = append(newConfig, result)
			} else if mode == "https" {
				result := ReplaceSSHtoHTTPS(scanner)

				newConfig = append(newConfig, line)
				newConfig = append(newConfig, result)
			}
		} else {
			newConfig = append(newConfig, line)
		}
	}

	config := strings.Join(newConfig, "\n")

	return config
}
