package commands

import (
	"context"
	"strings"

	"github.com/oam-dev/kubevela/api/types"
	"github.com/oam-dev/kubevela/pkg/oam"

	"github.com/gosuri/uitable"
	cmdutil "github.com/oam-dev/kubevela/pkg/commands/util"
	"github.com/spf13/cobra"
)

func NewTraitsCommand(c types.Args, ioStreams cmdutil.IOStreams) *cobra.Command {
	var workloadName string
	var syncCluster bool
	ctx := context.Background()
	cmd := &cobra.Command{
		Use:                   "traits [--apply-to WORKLOAD_NAME]",
		DisableFlagsInUseLine: true,
		Short:                 "List traits",
		Long:                  "List traits",
		Example:               `vela traits`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if syncCluster {
				if err := RefreshDefinitions(ctx, c, ioStreams, true); err != nil {
					return err
				}
			}
			return printTraitList(&workloadName, ioStreams)
		},
		Annotations: map[string]string{
			types.TagCommandType: types.TypeCap,
		},
	}

	cmd.SetOut(ioStreams.Out)
	cmd.Flags().StringVar(&workloadName, "apply-to", "", "Workload name")
	cmd.Flags().BoolVarP(&syncCluster, "sync", "s", true, "Synchronize capabilities from cluster into local")
	return cmd
}

func printTraitList(workloadName *string, ioStreams cmdutil.IOStreams) error {
	table := uitable.New()
	table.Wrap = true
	traitDefinitionList, err := oam.ListTraitDefinitions(workloadName)
	if err != nil {
		return err
	}
	table.AddRow("NAME", "DESCRIPTION", "APPLIES TO")
	simplifiedTraits := oam.SimplifyCapabilityStruct(traitDefinitionList)
	for _, t := range simplifiedTraits {
		table.AddRow(t.Name, t.Description, strings.Join(t.AppliesTo, "\n"))
	}
	ioStreams.Info(table.String())
	return nil
}
