package helm

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/helm"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
)

// DevSpaceChartConfig is the config that holds the devspace chart information
var DevSpaceChartConfig = &latest.ChartConfig{
	Name:    "component-chart",
	Version: "v0.0.6",
	RepoURL: "https://charts.devspace.cloud",
}

// DeployConfig holds the information necessary to deploy via helm
type DeployConfig struct {
	// Public because we can switch them to fake clients for testing
	Kube *kubectl.Client
	Helm helm.Interface

	TillerNamespace  string
	DeploymentConfig *latest.DeploymentConfig
	Log              log.Logger

	config *latest.Config
}

// New creates a new helm deployment client
func New(config *latest.Config, kubeClient *kubectl.Client, deployConfig *latest.DeploymentConfig, log log.Logger) (*DeployConfig, error) {
	tillerNamespace := kubeClient.Namespace
	if deployConfig.Helm.TillerNamespace != "" {
		tillerNamespace = deployConfig.Helm.TillerNamespace
	}

	// Exchange chart
	if deployConfig.Helm.ComponentChart != nil && *deployConfig.Helm.ComponentChart == true {
		deployConfig.Helm.Chart = DevSpaceChartConfig
	}

	return &DeployConfig{
		Kube:             kubeClient,
		TillerNamespace:  tillerNamespace,
		DeploymentConfig: deployConfig,
		Log:              log,
		config:           config,
	}, nil
}

// Delete deletes the release
func (d *DeployConfig) Delete(cache *generated.CacheConfig) error {
	// Delete with helm engine
	isDeployed := helm.IsTillerDeployed(d.config, d.Kube, d.TillerNamespace)
	if isDeployed == false {
		return nil
	}

	if d.Helm == nil {
		var err error

		// Get HelmClient
		d.Helm, err = helm.NewClient(d.config, d.Kube, d.TillerNamespace, d.Log, false)
		if err != nil {
			return errors.Wrap(err, "new helm client")
		}
	}

	_, err := d.Helm.DeleteRelease(d.DeploymentConfig.Name, true)
	if err != nil {
		return err
	}

	// Delete from cache
	delete(cache.Deployments, d.DeploymentConfig.Helm.Chart.Name)
	return nil
}
