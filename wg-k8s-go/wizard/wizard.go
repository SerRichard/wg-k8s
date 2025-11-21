package wizard

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Wizard struct{}

func (wizard Wizard) ValidateArgs(args []string) (bool, error) {

	k8sConfig, version, namespace, valuesFile := args[0], args[1], args[2], args[3]

	fmt.Printf("Running Wizard with: %v %v %v %v\n", k8sConfig, version, namespace, valuesFile)

	return true, nil
}


func (wizard Wizard) InitHelm() {
	
}

func (wizard Wizard) AddHelm(){
	
}

func (wizard Wizard) UpdateHelm() {

}

func (wizard Wizard) RunInstall(args []string) {
	
	fmt.Print("You ran the WgK8s Wizard!\n")

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


