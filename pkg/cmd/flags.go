package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func SupportUpdateSlices(cmd *cobra.Command, sliceName *string, sliceDefaultTag *string) {
	SupportBuildSlices(cmd, sliceName, sliceDefaultTag)
	cmd.Flags().StringVar(sliceDefaultTag, "slice-default-tag", "", "For stacks using slices, override the default tag for any images that aren't being built & pushed by the slice")
}

func SupportBuildSlices(cmd *cobra.Command, sliceName *string, sliceDefaultTag *string) {
	cmd.Flags().StringVarP(sliceName, "slice", "s", "", "If you only need to test a slice of the app, specify it here")
}

func ValidateUpdateSliceFlags(cmd *cobra.Command, args []string) error {
	if !(cmd.Flags().Changed("slice") == cmd.Flags().Changed("slice-default-tag")) {
		return errors.New("both (or None) of `slice` and `slice-default-tag` must be set (or unset).")
	}
	return nil
}
