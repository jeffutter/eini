package cmd

import (
	"fmt"
	"github.com/jeffutter/eini/crypto"
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		pubkey, err := cfg.PubKey()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		encrypter, err := crypto.PrepareEncrypter(pubkey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up Crypto: %s\n", err)
			os.Exit(1)
		}

		for _, sec := range cfg.GetSections() {
			for _, key := range sec.GetKeys() {
				if shouldEncrypt(key) {
					encrypted, err := encrypter.Encrypt([]byte(key.Value()))
					if err != nil {
						fmt.Fprintf(os.Stderr, "Failed encrypting key: %s\n", key.Name())
						os.Exit(1)
					}
					key.SetValue(fmt.Sprintf("%s", encrypted))
				}
			}
		}

		var output io.Writer
		if write {
			f, err := os.Create(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to open file %s for writing: %s\n", file, err)
				os.Exit(1)
			}
			output = f
			defer f.Close()
		} else {
			output = os.Stdout
		}

		_, err = cfg.WriteTo(output)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config: %s\n", err)
			os.Exit(1)
		}
	},
}

func shouldEncrypt(key ini.Key) bool {
	return !ignoreKeyRegex.MatchString(key.Name()) && !decryptedRegex.MatchString(key.Comment()) && !encryptedRegex.MatchString(key.Value())
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVarP(&write, "write", "w", false, "write back to original file")
}
