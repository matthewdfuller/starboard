package cmd

import (
	"errors"
	"strings"

	"github.com/aquasecurity/starboard/pkg/kube"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewRootCmd(version VersionInfo) *cobra.Command {
	var cf *genericclioptions.ConfigFlags

	rootCmd := &cobra.Command{
		Use:           "starboard",
		Short:         "Kubernetes-native security",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cf = genericclioptions.NewConfigFlags(true)

	rootCmd.AddCommand(NewVersionCmd(version))
	rootCmd.AddCommand(NewInitCmd(cf))
	rootCmd.AddCommand(NewRBACCmd(cf))
	rootCmd.AddCommand(NewFindCmd(cf))
	rootCmd.AddCommand(NewKubeBenchCmd(cf))
	rootCmd.AddCommand(NewKubeHunterCmd(cf))
	rootCmd.AddCommand(NewPolarisCmd(cf))
	rootCmd.AddCommand(NewGetCmd(cf))
	rootCmd.AddCommand(NewCleanupCmd(cf))

	SetFlags(cf, rootCmd)

	return rootCmd
}

func SetFlags(cf *genericclioptions.ConfigFlags, cmd *cobra.Command) {
	cf.AddFlags(cmd.Flags())
	for _, c := range cmd.Commands() {
		SetFlags(cf, c)
	}
}

func WorkloadFromArgs(namespace string, args []string) (workload kube.Workload, err error) {
	if len(args) < 1 {
		err = errors.New("required workload kind and name not specified")
		return
	}

	parts := strings.SplitN(args[0], "/", 2)
	if len(parts) == 1 {
		workload = kube.Workload{
			Namespace: namespace,
			Kind:      kube.WorkloadKindPod,
			Name:      parts[0],
		}
		return
	}
	kind, err := kube.WorkloadKindFromString(parts[0])
	if err != nil {
		return
	}
	if "" == parts[1] {
		err = errors.New("required workload name is blank")
		return
	}
	workload = kube.Workload{
		Namespace: namespace,
		Kind:      kind,
		Name:      parts[1],
	}
	return
}
