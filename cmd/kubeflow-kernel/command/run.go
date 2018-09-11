package command

import (
	"log"

	kubeflowbackend "github.com/caicloud/ciao/pkg/backend/kubeflow"
	simpleinterpreter "github.com/caicloud/ciao/pkg/interpreter/simple"
	"github.com/caicloud/ciao/pkg/kernel"
	"github.com/caicloud/ciao/pkg/manager"
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
	kubeConfig := viper.GetString("kubeconfig")
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

	s2iClient := simples2i.New()

	interpreter := simpleinterpreter.New()

	mgr := manager.New(kubeflowBackend, s2iClient, interpreter)

	ciao := kernel.New(version.ProtocolVersion, version.Version, connectionFile, mgr)

	log.Println("Running Kubeflow kernel for Jupyter...")
	ciao.RunKernel()
}
