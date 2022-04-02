package cmd

import (
	"errors"
	"fmt"
	"os"
	"vuerd/engines/ent"
	"vuerd/engines/prisma"
	"vuerd/types"
	"vuerd/utils"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vuerd",
	Short: "Generate API based on your ERD diagram ( vuerd vscode extension )",
	Run: func(cmd *cobra.Command, args []string) {
		vuerdCmd()
	},
}

type promptContent struct {
	errorMsg string
	label    string
}

func vuerdCmd() {
	dbPath := promptGetInput(promptContent{
		label:    "Enter the path of your database ",
		errorMsg: "Please enter the path of your database ",
	})

	pkg := promptGetInput(promptContent{
		label:    "Enter the package name of your project",
		errorMsg: "Please enter the package name of your project",
	})

	// dbType := promptGetSelect(promptContent{
	// 	label: "Select the type of your database",
	// }, []string{"mysql", "postgres", "sqlite"})

	schema := promptGetSelect(promptContent{
		label: "Select ORM:",
	}, []string{"ent", "prisma"})

	var state types.State
	utils.ReadJSON(&state, dbPath)

	switch schema {
	case "ent":
		ent.Ent(state, pkg)
	case "prisma":
		prisma.Prisma(state)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}

		return nil
	}
	template := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: template,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Input: %s\n", result)
	return result
}

func promptGetSelect(pc promptContent, options []string) string {
	prompt := promptui.SelectWithAdd{
		Label: pc.label,
		Items: options,
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}
