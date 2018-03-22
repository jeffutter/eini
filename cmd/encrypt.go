package cmd

import (
	"github.com/jeffutter/eini/ini"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
)

var write bool

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypts an eini file",
	Long: `(re-)encrypt an eini file using the
public key contained in the _public_key entry.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Encryption is expensive. We'd rather burn cycles on many cores than wait.
		runtime.GOMAXPROCS(runtime.NumCPU())

		file := args[0]

		cfg, err := ini.Load(file)
		checkError(err)

		pubkey, err := cfg.PubKey()
		checkError(err)

		var output io.Writer
		if write {
			f, err := os.Create(file)
			checkErrorf(err, "Unable to open file %s for writing: %s\n", file, err)

			output = f
			defer f.Close()
		} else {
			output = os.Stdout
		}

		err = cfg.Encrypt(pubkey, output)
		checkErrorf(err, "Failed to write config: %s\n", err)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVarP(&write, "write", "w", false, "write back to original file")
}
