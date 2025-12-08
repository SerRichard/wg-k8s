package wizard

import (
    "context"
    "fmt"
    "log"
    "os"
    "encoding/base64"

	"github.com/spf13/cobra"
    "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

    coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
    k8serrors "k8s.io/apimachinery/pkg/api/errors"
	
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/chart/loader"
    "helm.sh/helm/v3/pkg/cli"
    "helm.sh/helm/v3/pkg/cli/values"
    "helm.sh/helm/v3/pkg/getter"
    "helm.sh/helm/v3/pkg/repo"
)

type Wizard struct{}
type Keys struct{
    PrivateKey wgtypes.Key
    PublicKey wgtypes.Key
}

func (wizard Wizard) ValidateArgs(args []string) (bool, error) {

	k8sConfig, version, namespace, valuesFile := args[0], args[1], args[2], args[3]

	fmt.Printf("Running Wizard with: %v\n %v\n %v\n %v\n", k8sConfig, version, namespace, valuesFile)

	return true, nil
}

func (wizard Wizard) KubernetesClient(kubeconfig string) coreV1Types.CoreV1Interface  {

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset.CoreV1()
}

func (wizard Wizard) GenerateKeys() (Keys) {
    privateKey, err := wgtypes.GeneratePrivateKey()
    if err != nil {
        panic(err)
    }

    keys := Keys{}
    keys.PrivateKey = privateKey
    keys.PublicKey = privateKey.PublicKey()

    return keys
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

    install := action.NewUpgrade(actionConfig)
    install.Namespace = namespace
    // install.ReleaseName = releaseName
	install.ChartPathOptions.Version = version
    install.Install = true

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
    _, err = install.Run(releaseName, chart, vals)
    return err
}

func (wizard Wizard) RunInstall(args []string) {

	k8sConfig, version, namespace, valuesFile := args[0], args[1], args[2], args[3]

    // Prepare keys
    keys := wizard.GenerateKeys()

    K8sclient := wizard.KubernetesClient(k8sConfig)

    // Prepare namespace
    nsName := &coreV1.Namespace{
        ObjectMeta: metaV1.ObjectMeta{
            Name: namespace,
        },
    }
    
    _, err := K8sclient.Namespaces().Create(context.Background(), nsName, metaV1.CreateOptions{})
    
    if err != nil {
        if !k8serrors.IsAlreadyExists(err) {
            panic(err)
        }
    }

	// Prepare secret
    secretsClient := K8sclient.Secrets(namespace)

    privateKeyB64 := base64.StdEncoding.EncodeToString(keys.PrivateKey[:])
    secretSpec := &coreV1.Secret{
        ObjectMeta: metaV1.ObjectMeta{
            Name: "wireguard-secret",
            Namespace: namespace,
        },
        Data: map[string][]byte{
            "privatekey": []byte(privateKeyB64),
        },
    }

	_, err = secretsClient.Create(
        context.Background(),
        secretSpec,
        metaV1.CreateOptions{},
    )
	if err != nil {
        if k8serrors.IsAlreadyExists(err) {

            existing, getErr := secretsClient.Get(context.Background(), secretSpec.Name, metaV1.GetOptions{})
            if getErr != nil {
                panic(getErr)
            }

            secretSpec.ResourceVersion = existing.ResourceVersion

            _, updateErr := secretsClient.Update(context.Background(), secretSpec, metaV1.UpdateOptions{})
            if updateErr != nil {
                fmt.Printf("Two")
                panic(updateErr)
            } 
            
        } else {
                fmt.Printf("Three %v", err)
                panic(err)
        }
    }  else {
            fmt.Println("Created new secret.")
    }
	// Install chart!
	err = wizard.InstallChart(    
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


