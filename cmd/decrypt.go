package cmd

import (
	"fmt"
	"github.com/jeffutter/eini/crypto"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var output string
var keydir string

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
		ignoreKeyRegex, _ := regexp.Compile("^_.*")

		cfg, err := ini.Load(args[0])
		if err != nil {
			fmt.Printf("Fail to read file %s: %v", args[0], err)
			return
		}

		pubkey, err := cfg.Section("").GetKey("_public_key")
		if err != nil {
			fmt.Printf("Couldn't read public key from ini")
			return
		}

		stdinContent, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read from stdin:", err)
			os.Exit(1)
		}

		privKey := strings.TrimSpace(string(stdinContent))

		if privKey == "" {
			privKey, err = readPrivateKeyFromDisk(pubkey.Value(), keydir)

			if err != nil {
				fmt.Printf("Error reading private key from disk: %s", err)
				return
			}

			if privKey == "" {
				fmt.Printf("Private key not provided, aborting")
				return
			}
		}

		decrypter, err := crypto.PrepareDecrypter(pubkey.Value(), privKey)
		if err != nil {
			fmt.Printf("Error setting up Crypto: %s\n", err)
			return
		}

		for _, sec := range cfg.SectionStrings() {
			section, err := cfg.GetSection(sec)
			if err != nil {
				fmt.Printf("Failed parsing ini section %s\n", sec)
				return
			}
			for _, key := range section.KeyStrings() {
				if !ignoreKeyRegex.MatchString(key) {
					val := section.Key(key).Value()
					decrypted, err := crypto.Decrypt(decrypter, val)
					if err != nil {
						fmt.Printf("Failed decrypting key: %s\n", sec)
						return
					}
					if sec == "DEFAULT" {
						fmt.Printf("declare -x \"%s\"=\"%s\"\n", key, decrypted)
					} else {
						fmt.Printf("declare -x \"%s_%s\"=\"%s\"\n", strings.ToUpper(sec), key, decrypted)
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
	decryptCmd.Flags().StringVarP(&output, "output", "o", "env", "output format: [env, yaml]")
	decryptCmd.Flags().StringVarP(&keydir, "keydir", "k", "/opt/ejson/keys", "Directory containing EJSON keys")
}
