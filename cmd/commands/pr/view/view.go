package view

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View a pull request",
		Long:  "Display the title, body, and other information about a pull request.",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			c := internal.Client{
				Username: globalOpts.Username,
				Password: globalOpts.Password,
			}

			bbrepo, err := bbgit.GetBitbucketRepo()

			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			if !bbrepo.IsBitbucketOrg() {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Are you sure this is a bitbucket repo?")
				return
			}

			if len(args) > 0 {
				id, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
				if len(prs.Values) == 0 {
					fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
					return
				}

				id = prs.Values[0].ID

			}

			pr, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			if Web {
				err := browser.OpenURL(pr.Links["html"].Href)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
				return
			}

			commits, err := c.PrCommits(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			PrintSummary(pr, commits)
		},
	}
	viewCmd.Flags().BoolVar(&Web, "web", false, "view the pull request in your browser")
	prCmd.AddCommand(viewCmd)
}

func PrintSummary(pr *internal.PullRequest, commits *internal.Commits) {
	fmt.Println(aurora.Bold(pr.Title))
	var state string
	if pr.State == "OPEN" {
		state = aurora.Green("Open").String()
	} else if pr.State == "DECLINED" {
		state = aurora.Red("Declined").String()
	} else {
		state = pr.State
	}

	nrOfCommits := len(commits.Values)
	nrOfCommitsStr := fmt.Sprintf("%d", nrOfCommits)
	if nrOfCommits == 10 {
		nrOfCommitsStr = nrOfCommitsStr + "+"
	}

	infoText := aurora.Index(242, fmt.Sprintf("%s wants to merge %s commits into %s from %s\n", pr.Author.Nickname, nrOfCommitsStr, pr.Destination.Branch.Name, pr.Source.Branch.Name))
	fmt.Printf("%s • %s\n", state, infoText)
	if pr.Description != "" {
		out, err := glamour.Render(pr.Description, "dark")
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			return
		}
		fmt.Println(out)
	}

	footer := aurora.Index(242, fmt.Sprintf("View this pull request on Bitbucket.org: %s", pr.Links["html"].Href)).String()
	fmt.Println(footer)
	// fmt.Println(pr, err)

}
