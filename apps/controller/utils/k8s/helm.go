package k8s

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func newHelmActionCfg(namespace string) (*action.Configuration, error) {
	actionCfg := new(action.Configuration)
	helmEnv := cli.New()

	if err := actionCfg.Init(helmEnv.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	return actionCfg, nil
}

func helmInstall(actionCfg *action.Configuration, namespace string, releaseName string, chartPath string, chartValues map[string]any) (*release.Release, error) {
	install := action.NewInstall(actionCfg)
	install.Namespace = namespace
	install.ReleaseName = releaseName
	install.CreateNamespace = true
	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	release, err := install.Run(chart, chartValues)
	return release, err
}

func helmUninstall(actionCfg *action.Configuration, releaseName string) (*release.UninstallReleaseResponse, error) {
	uninstall := action.NewUninstall(actionCfg)
	res, err := uninstall.Run(releaseName)
	return res, err
}
