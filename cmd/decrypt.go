package cmd

import (
	"fmt"
	"github.com/jeffutter/eini/crypto"
	"github.com/jeffutter/eini/ini"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var output string
var keydir string

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts an eini file",
	Long: `Decrypts an eini file and prints it to stdout using the
private key passed to stdin or in the keydir.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := ini.Load(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		pubkey, err := cfg.PubKey()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		stdinContent, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read from stdin:", err)
			os.Exit(1)
		}

		privKey := strings.TrimSpace(string(stdinContent))

		if privKey == "" {
			privKey, err = readPrivateKeyFromDisk(pubkey, keydir)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading private key from disk: %s", err)
				os.Exit(1)
			}

			if privKey == "" {
				fmt.Fprintf(os.Stderr, "Private key not provided, aborting")
				os.Exit(1)
			}
		}

		decrypter, err := crypto.PrepareDecrypter(pubkey, privKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up Crypto: %s\n", err)
			os.Exit(1)
		}

		for _, sec := range cfg.GetSections() {
			keys := sec.GetKeys()
			for _, key := range keys {
				if !ignoreKeyRegex.MatchString(key.Name()) {
					decrypted, err := crypto.Decrypt(decrypter, key.Value())
					if err != nil {
						fmt.Fprintf(os.Stderr, "Failed decrypting key: %s\n", sec.Name())
						os.Exit(1)
					}
					if sec.Name() == "DEFAULT" {
						fmt.Printf("declare -x \"%s\"=\"%s\"\n", key.Name(), decrypted)
					} else {
						fmt.Printf("declare -x \"%s_%s\"=\"%s\"\n", strings.ToUpper(sec.Name()), key.Name(), decrypted)
					}
				}
			}
		}
	},
}

func readPrivateKeyFromDisk(pubkey string, keydir string) (privkey string, err error) {
	keyFile := fmt.Sprintf("%s/%s", keydir, pubkey)
	var fileContents []byte
	fileContents, err = ioutil.ReadFile(keyFile)
	if err != nil {
		err = fmt.Errorf("couldn't read key file (%s)", err.Error())
		return
	}
	privkey = string(fileContents)
	return privkey, nil
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&output, "output", "o", "env", "output format: [env]")
	decryptCmd.Flags().StringVarP(&keydir, "keydir", "k", "/opt/ejson/keys", "Directory containing EJSON keys")
}
