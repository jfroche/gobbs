package pullrequests

import (
	"net/http"
	"strconv"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	"github.com/spf13/cobra"
)

var (
	prVersion *int64
)

func init() {
	prVersion = Delete.Flags().Int64P("version", "V", 0, "Define version of PR to delete (Modified PR's increase version)")
}

// Delete is the cmd implementation for Deleting Pull Requests
var Delete = &cobra.Command{
	Use:     "delete pullRequestID",
	Aliases: []string{"del"},
	Short:   "Delete pull requests for repository",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Argument must be a pull request ID. Err: %s", err.Error())
		}

		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Fatal(err.Error())
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)

		response, err := apiClient.DefaultApi.DeleteWithVersion(*stashInfo.Project(), *stashInfo.Repo(), prID, *prVersion)

		if response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		if response.StatusCode == http.StatusNoContent {
			log.Infof("Pull request ID: %v sucessfully deleted.", prID)
		}
	},
}