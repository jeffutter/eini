package cmd

import (
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
		checkError(err)

		pubkey, err := cfg.PubKey()
		checkError(err)

		encrypter, err := crypto.PrepareEncrypter(pubkey)
		checkErrorf(err, "Error setting up Crypto: %s\n", err)

		for _, sec := range cfg.GetSections() {
			for _, key := range sec.GetKeys() {
				if shouldEncrypt(key) {
					encrypted, err := encrypter.Encrypt(key.Value())
					checkErrorf(err, "Failed encrypting key: %s\n", key.Name())

					key.SetValue(encrypted)
				}
			}
		}

		var output io.Writer
		if write {
			f, err := os.Create(file)
			checkErrorf(err, "Unable to open file %s for writing: %s\n", file, err)

			output = f
			defer f.Close()
		} else {
			output = os.Stdout
		}

		_, err = cfg.WriteTo(output)
		checkErrorf(err, "Failed to write config: %s\n", err)
	},
}

func shouldEncrypt(key ini.Key) bool {
	return !ignoreKeyRegex.MatchString(key.Name()) && !decryptedRegex.MatchString(key.Comment()) && !encryptedRegex.MatchString(key.Value())
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVarP(&write, "write", "w", false, "write back to original file")
}
