package wizard

import (
    "fmt"
    "log"
    "os"

	"github.com/spf13/cobra"
	
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/chart/loader"
    "helm.sh/helm/v3/pkg/cli"
    "helm.sh/helm/v3/pkg/cli/values"
    "helm.sh/helm/v3/pkg/getter"
    "helm.sh/helm/v3/pkg/repo"
)

type Wizard struct{}

func (wizard Wizard) ValidateArgs(args []string) (bool, error) {

	k8sConfig, version, namespace, valuesFile := args[0], args[1], args[2], args[3]

	fmt.Printf("Running Wizard with: %v\n %v\n %v\n %v\n", k8sConfig, version, namespace, valuesFile)

	return true, nil
}

func (wizard Wizard) InstallChart(kubeconfig, namespace, releaseName, repoName, repoURL, chartName, version, valuesFile string) error {
    settings := cli.New()
    settings.KubeConfig = kubeconfig
    settings.SetNamespace(namespace)

    // Load or init repositories.yaml
    rf, err := repo.LoadFile(settings.RepositoryConfig)
    if err != nil && !os.IsNotExist(err) {
        return err
    }
    if rf == nil {
        rf = repo.NewFile()
    }

    // Add repo if missing
    entry := &repo.Entry{Name: repoName, URL: repoURL}
    if !rf.Has(entry.Name) {
        chartRepo, err := repo.NewChartRepository(entry, getter.All(settings))
        if err != nil {
            return err
        }
        if _, err := chartRepo.DownloadIndexFile(); err != nil {
            return err
        }
        rf.Update(entry)
        if err := rf.WriteFile(settings.RepositoryConfig, 0644); err != nil {
            return err
        }
    }

    // Update repo
	r, err := repo.NewChartRepository(entry, getter.All(settings))
	if err != nil {
		return err
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return err
	}

    // Action config
    actionConfig := new(action.Configuration)
    if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
        return err
    }

    install := action.NewInstall(actionConfig)
    install.Namespace = namespace
    install.ReleaseName = releaseName
	install.ChartPathOptions.Version = version

    // Locate & download chart
    repoChart := fmt.Sprintf("%s/%s", repoName, chartName)
    chartPath, err := install.ChartPathOptions.LocateChart(repoChart, settings)
    if err != nil {
        return err
    }

    // Load chart
    chart, err := loader.Load(chartPath)
    if err != nil {
        return err
    }

    // Load values
    valOpts := &values.Options{ValueFiles: []string{valuesFile}}
    vals, err := valOpts.MergeValues(getter.All(settings))
    if err != nil {
        return err
    }

    // Install
    _, err = install.Run(chart, vals)
    return err
}

func (wizard Wizard) RunInstall(args []string) {

	k8sConfig, version, namespace, valuesFile := args[0], args[1], args[2], args[3]

	// Install wg CLI!
	
	// Prepare secret

	// Prepare namespace

	// Install chart!
	err := wizard.InstallChart(    
		k8sConfig,
		namespace,
		"wg-k8s",
		"wg-k8s",
		"https://serrichard.github.io/wg-k8s",
		"wg-k8s",
		version,
		valuesFile,
	)
	if err != nil {
		fmt.Printf("Failed to install %v\n", err)
	}
		
}

func WizardCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "wizard <k8s-config> <version> <namespace> <values>",
		Short: "WgK8s Wizard!",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			wizard := Wizard{}

			wizard.ValidateArgs(args)

			wizard.RunInstall(args)
		},
	}

	command.Flags().Bool("output", true, "output peer config")

	return command	
}


