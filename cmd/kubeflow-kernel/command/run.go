package command

import (
	"fmt"
	"log"

	kubeflowbackend "github.com/caicloud/ciao/pkg/backend/kubeflow"
	"github.com/caicloud/ciao/pkg/config"
	simpleinterpreter "github.com/caicloud/ciao/pkg/interpreter/simple"
	"github.com/caicloud/ciao/pkg/kernel"
	"github.com/caicloud/ciao/pkg/manager"
	"github.com/caicloud/ciao/pkg/s2i"
	imgs2i "github.com/caicloud/ciao/pkg/s2i/img"
	simples2i "github.com/caicloud/ciao/pkg/s2i/simple"
	"github.com/caicloud/ciao/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
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
	if kubeConfig == "" {
		log.Fatalln("Failed to start the kernel: Kubeconfig missed")
	}

	// Get kubernetes config.
	kcfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Error building kubeConfig: %s\n", err.Error())
	}

	kubeflowBackend, err := kubeflowbackend.New(kcfg)
	if err != nil {
		log.Fatalf("Error building kubeflow backend: %s\n", err.Error())
	}

	s2iConfig := viper.GetStringMapString(config.S2I)
	if s2iConfig == nil {
		log.Fatalf("Error creating s2i client: Failed to find the config\n")
	}

	s2iClient, err := createS2IClient(s2iConfig)
	if err != nil {
		log.Fatalf("Error creating s2i client: %s\n", err.Error())
	}

	interpreter := simpleinterpreter.New()

	mgr := manager.New(kubeflowBackend, s2iClient, interpreter)

	ciao := kernel.New(version.ProtocolVersion, version.Version, connectionFile, mgr)

	log.Println("Running Kubeflow kernel for Jupyter...")
	ciao.RunKernel()
}

func createS2IClient(s2iConfig map[string]string) (s2i.Interface, error) {
	switch s2iConfig[config.S2IProvider] {
	case config.S2IProviderS2I:
		return simples2i.New(), nil
	case config.S2IProviderImg:
		return imgs2i.New(s2iConfig[config.S2IRegistry], s2iConfig[config.S2IUsername], s2iConfig[config.S2IPassword])
	default:
		return nil, fmt.Errorf("Failed to find the provider %s", s2iConfig[config.S2IProvider])
	}
}
