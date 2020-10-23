package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/werf/werf/cmd/werf/synchronization"

	"github.com/werf/werf/cmd/werf/slugify"

	"github.com/werf/werf/cmd/werf/build"
	"github.com/werf/werf/cmd/werf/ci_env"
	"github.com/werf/werf/cmd/werf/run"

	"github.com/werf/werf/cmd/werf/cleanup"
	"github.com/werf/werf/cmd/werf/purge"

	"github.com/werf/werf/cmd/werf/dismiss"

	"github.com/sirupsen/logrus"

	"github.com/werf/logboek"

	"github.com/spf13/cobra"

	"github.com/werf/werf/cmd/werf/converge"
	"github.com/werf/werf/cmd/werf/helm"

	managed_images_add "github.com/werf/werf/cmd/werf/managed_images/add"
	managed_images_ls "github.com/werf/werf/cmd/werf/managed_images/ls"
	managed_images_rm "github.com/werf/werf/cmd/werf/managed_images/rm"

	host_cleanup "github.com/werf/werf/cmd/werf/host/cleanup"
	host_project_list "github.com/werf/werf/cmd/werf/host/project/list"
	host_project_purge "github.com/werf/werf/cmd/werf/host/project/purge"
	host_purge "github.com/werf/werf/cmd/werf/host/purge"

	config_list "github.com/werf/werf/cmd/werf/config/list"
	config_render "github.com/werf/werf/cmd/werf/config/render"

	"github.com/werf/werf/cmd/werf/completion"
	"github.com/werf/werf/cmd/werf/docs"
	"github.com/werf/werf/cmd/werf/version"

	"github.com/werf/werf/cmd/werf/common"
	"github.com/werf/werf/cmd/werf/common/templates"
	"github.com/werf/werf/pkg/process_exterminator"
)

func main() {
	common.EnableTerminationSignalsTrap()
	log.SetOutput(logboek.ProxyOutStream())
	logrus.StandardLogger().SetOutput(logboek.ProxyOutStream())

	if err := process_exterminator.Init(); err != nil {
		common.TerminateWithError(fmt.Sprintf("process exterminator initialization failed: %s", err), 1)
	}

	rootCmd := constructRootCmd()

	if err := rootCmd.Execute(); err != nil {
		common.TerminateWithError(err.Error(), 1)
	}
}

func constructRootCmd() *cobra.Command {
	if filepath.Base(os.Args[0]) == "helm" || os.Getenv("WERF_HELM3_MODE") == "1" {
		return helm.NewCmd()
	}

	rootCmd := &cobra.Command{
		Use:   "werf",
		Short: "werf helps to implement and support Continuous Integration and Continuous Delivery",
		Long: common.GetLongCommandDescription(`werf helps to implement and support Continuous Integration and Continuous Delivery.

Find more information at https://werf.io`),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	groups := &templates.CommandGroups{}
	*groups = append(*groups, templates.CommandGroups{
		{
			Message: "Delivery commands",
			Commands: []*cobra.Command{
				converge.NewCmd(),
				dismiss.NewCmd(),
			},
		},
		{
			Message: "Cleaning commands",
			Commands: []*cobra.Command{
				cleanup.NewCmd(),
				purge.NewCmd(),
			},
		},
		{
			Message: "Helper commands",
			Commands: []*cobra.Command{
				ci_env.NewCmd(),
				build.NewCmd(),
				run.NewCmd(),
				slugify.NewCmd(),
			},
		},
		{
			Message: "Low-level management commands",
			Commands: []*cobra.Command{
				configCmd(),
				managedImagesCmd(),
				hostCmd(),
				helm.NewCmd(),
			},
		},
		{
			Message: "Other commands",
			Commands: []*cobra.Command{
				synchronization.NewCmd(),
				completion.NewCmd(rootCmd),
				version.NewCmd(),
				docs.NewCmd(groups),
			},
		},
	}...)
	groups.Add(rootCmd)

	templates.ActsAsRootCommand(rootCmd, *groups...)

	return rootCmd
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Work with werf.yaml",
	}
	cmd.AddCommand(
		config_render.NewCmd(),
		config_list.NewCmd(),
	)

	return cmd
}

func managedImagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "managed-images",
		Short: "Work with managed images which will be preserved during cleanup procedure",
	}
	cmd.AddCommand(
		managed_images_add.NewCmd(),
		managed_images_ls.NewCmd(),
		managed_images_rm.NewCmd(),
	)

	return cmd
}

func hostCmd() *cobra.Command {
	hostCmd := &cobra.Command{
		Use:   "host",
		Short: "Work with werf cache and data of all projects on the host machine",
	}

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Work with projects",
	}

	projectCmd.AddCommand(
		host_project_list.NewCmd(),
		host_project_purge.NewCmd(),
	)

	hostCmd.AddCommand(
		host_cleanup.NewCmd(),
		host_purge.NewCmd(),
		projectCmd,
	)

	return hostCmd
}
