package helm

import (
	"context"
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/werf/werf/pkg/werf"

	"github.com/werf/werf/pkg/deploy/werf_chart"

	"github.com/spf13/cobra"
	cmd_werf_common "github.com/werf/werf/cmd/werf/common"
	cmd_helm "helm.sh/helm/v3/cmd/helm"
	"helm.sh/helm/v3/pkg/action"
)

var upgradeCmdData cmd_werf_common.CmdData

func NewUpgradeCmd(actionConfig *action.Configuration) *cobra.Command {
	wc := werf_chart.NewWerfChart(werf_chart.WerfChartOptions{})

	cmd, helmAction := cmd_helm.NewUpgradeCmd(actionConfig, os.Stdout, cmd_helm.UpgradeCmdOptions{
		LoadOptions: loader.LoadOptions{
			ChartExtender:               wc,
			SubchartExtenderFactoryFunc: func() chart.ChartExtender { return werf_chart.NewWerfChart(werf_chart.WerfChartOptions{}) },
		},
		PostRenderer: wc.ExtraAnnotationsAndLabelsPostRenderer,
	})

	SetupWerfChartParams(cmd, &upgradeCmdData)

	oldRunE := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := werf.Init(*upgradeCmdData.TmpDir, *upgradeCmdData.HomeDir); err != nil {
			return err
		}

		if chartDir, err := helmAction.ChartPathOptions.LocateChart(args[1], cmd_helm.Settings); err != nil {
			return err
		} else {
			wc.ChartDir = chartDir
		}
		wc.ReleaseName = args[0]

		if err := InitWerfChartParams(&upgradeCmdData, wc, wc.ChartDir); err != nil {
			return fmt.Errorf("unable to init werf chart: %s", err)
		}

		return wc.WrapUpgrade(context.Background(), func() error {
			return oldRunE(cmd, args)
		})
	}

	return cmd
}
