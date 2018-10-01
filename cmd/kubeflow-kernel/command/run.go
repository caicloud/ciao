package command

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	restclientset "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/caicloud/ciao/pkg/backend"
	kubeflowbackend "github.com/caicloud/ciao/pkg/backend/kubeflow"
	"github.com/caicloud/ciao/pkg/config"
	simpleinterpreter "github.com/caicloud/ciao/pkg/interpreter/simple"
	"github.com/caicloud/ciao/pkg/kernel"
	"github.com/caicloud/ciao/pkg/manager"
	"github.com/caicloud/ciao/pkg/s2i"
	configs2i "github.com/caicloud/ciao/pkg/s2i/configmap"
	imgs2i "github.com/caicloud/ciao/pkg/s2i/img"
	simples2i "github.com/caicloud/ciao/pkg/s2i/simple"
	"github.com/caicloud/ciao/version"
)

var connectionFile string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the kernel",
	Long:  ``,
	Run:   run,
}

func init() {
	runCmd.Flags().StringVar(&connectionFile, "connection-file", "", "Connection File (which is assigned by Jupyter)")
	RootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) {
	kubeConfig := viper.GetString(config.KubeConfig)

	// Get kubernetes config.
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Error building kubeConfig: %s\n", err.Error())
	}

	s2iConfig := viper.GetStringMapString(config.S2I)
	if s2iConfig == nil {
		log.Fatalf("Error creating s2i client: Failed to find the config\n")
	}

	kubeflowBackend, err := createBackend(s2iConfig, kcfg)
	if err != nil {
		log.Fatalf("Error building kubeflow backend: %s\n", err.Error())
	}

	s2iClient, err := createS2IClient(s2iConfig, kcfg)
	if err != nil {
		log.Fatalf("Error creating s2i client: %s\n", err.Error())
	}

	interpreter := simpleinterpreter.New()

	mgr := manager.New(kubeflowBackend, s2iClient, interpreter)

	ciao := kernel.New(version.ProtocolVersion, version.Version, connectionFile, mgr)

	log.Println("Running Kubeflow kernel for Jupyter...")
	ciao.RunKernel()
}

func createS2IClient(s2iConfig map[string]string, kubeconfig *restclientset.Config) (s2i.Interface, error) {
	switch s2iConfig[config.S2IProvider] {
	case config.S2IProviderS2I:
		return simples2i.New(), nil
	case config.S2IProviderImg:
		return imgs2i.New(s2iConfig[config.S2IRegistry], s2iConfig[config.S2IUsername], s2iConfig[config.S2IPassword])
	case config.S2IProviderCM:
		return configs2i.New(kubeconfig)
	default:
		return nil, fmt.Errorf("Failed to find the provider %s", s2iConfig[config.S2IProvider])
	}
}

func createBackend(s2iConfig map[string]string, kubeconfig *restclientset.Config) (backend.Interface, error) {
	switch s2iConfig[config.S2IProvider] {
	case config.S2IProviderCM:
		return kubeflowbackend.NewWithCM(kubeconfig)
	default:
		return kubeflowbackend.New(kubeconfig)
	}
}
