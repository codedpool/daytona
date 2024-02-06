// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd_profile

import (
	"os"

	profile_list "github.com/daytonaio/daytona/cmd/views/profilie_list"
	"github.com/daytonaio/daytona/config"
	"github.com/daytonaio/daytona/internal/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var profileDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete profile [PROFILE_NAME]",
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"remove", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		var chosenProfileId string
		var chosenProfile config.Profile

		if len(args) == 0 {
			profilesList := c.Profiles

			chosenProfileId = profile_list.GetProfileIdFromPrompt(profilesList, c.ActiveProfileId, "Choose a profile to delete", false)

			if chosenProfileId == "" {
				return
			}
		} else {
			chosenProfileId = args[0]
		}

		if chosenProfileId == "default" {
			log.Fatal("Can not delete default profile")
		}

		for _, profile := range c.Profiles {
			if profile.Id == chosenProfileId || profile.Name == chosenProfileId {
				chosenProfile = profile
				break
			}
		}

		if chosenProfile == (config.Profile{}) {
			log.Fatal("Profile does not exist")
			return
		}

		if c.ActiveProfileId == chosenProfile.Id {
			c.ActiveProfileId = "default"
		}

		for _, profile := range c.Profiles {
			if profile.Name == chosenProfile.Name || profile.Id == chosenProfile.Id {
				err = c.RemoveProfile(profile.Id)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}

		err = os.RemoveAll(util.GetProfilePath(chosenProfile.Id))
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Deleted profile %s", chosenProfile.Name)
	},
}