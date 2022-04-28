/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/terakilobyte/onboarder/genssh"
	"github.com/terakilobyte/onboarder/ghops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
)

var outDir string
var team string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onboarder",
	Short: "Bootstrap your work git repositories.",
	Long: `Onboarder will bootstrap your work git repositories.
It will authorize to Github and use an oauth credential to create forks.
It will generate a new ssh keypair and upload the public key to Github.
It will clone your forks, and set upstreams appropriately for each one.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

IMPORTANT: You will be asked if a question during the process similar to:

  The authenticity of host 'github.com (140.82.112.4)' can't be established.
  ED25519 key fingerprint is SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU.
  This key is not known by any other names
  Are you sure you want to continue connecting (yes/no/[fingerprint])?

You must answer yes to this question. It is adding the fingerprint to your
known_hosts file.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
I'm about to begin forking and cloning all of the repositories that you should need.
I'm also going to create an ssh key for you and add it to your github account.

This may take a while (5-10 minutes) depending on how many repositories I'm
working with. Please be patient.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

IMPORTANT: You will be asked if a question during the process similar to:

  The authenticity of host 'github.com (140.82.112.4)' can't be established.
  ED25519 key fingerprint is SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU.
  This key is not known by any other names
  Are you sure you want to continue connecting (yes/no/[fingerprint])?

You must answer yes to this question.
`)
		client, user, err := ghops.InitClient(ghops.AuthToGithub())
		if err != nil {
			log.Fatalln(err)
		}

		pubkey, keypath, err := genssh.SetupSSH(*user)
		if err != nil {
			log.Fatalln(err)
		}
		ghops.UploadKey(client, pubkey)
		ghops.ForkRepos(client, globals.GetReposForTeam(team))

		fmt.Println("Waiting 30 seconds for forks to complete...")
		time.Sleep(30 * time.Second)

		gitops.SetupLocalRepos(globals.GetReposForTeam(team), *user, outDir, keypath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&outDir, "out-dir", "o", "", "output directory")
	rootCmd.Flags().StringVarP(&team, "team", "t", "", "team name")
	cobra.MarkFlagRequired(rootCmd.Flags(), "out-dir")
	cobra.MarkFlagRequired(rootCmd.Flags(), "team")
}
