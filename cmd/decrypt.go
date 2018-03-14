package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var output string

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Decrypting to: " + output)
		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Print(err)
		}
		var val interface{}
		err = yaml.Unmarshal(b, &val)
		if err != nil {
			fmt.Println("unmarshal []byte to yaml failed: " + err.Error())
		}
		fmt.Printf("%v\n", val)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	decryptCmd.Flags().StringVarP(&output, "output", "o", "env", "output format: [env, yaml]")
}
