package cmd

import (
	"fmt"
	"github.com/jeffutter/eini/crypto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// keygenCmd represents the encrypt command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		pub, priv, err := crypto.GenerateKeypair()
		if err != nil {
			fmt.Printf("Error generating Keypair: %s\n", err)
			return
		}

		if write {
			keyFile := fmt.Sprintf("%s/%s", keydir, pub)
			err := ioutil.WriteFile(keyFile, append([]byte(priv), '\n'), 0440)
			if err != nil {
				fmt.Printf("Error writing keyfile: %s\n", err)
				return
			}
			fmt.Println(pub)
		} else {
			fmt.Fprintf(os.Stderr, "Public Key:\n%s\nPrivate Key:\n%s\n", pub, priv)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
	keygenCmd.Flags().BoolVarP(&write, "write", "w", false, "writes the private key to the KEYDIR")
	keygenCmd.Flags().StringVarP(&keydir, "keydir", "k", "/opt/ejson/keys", "Directory containing EJSON keys")
}
