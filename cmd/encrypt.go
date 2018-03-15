package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/Shopify/ejson/crypto"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io"
	"os"
	"regexp"
)

var write bool

func prepareEncrypter(pubKey string) *crypto.Encrypter {
	var pub [32]byte

	pubkey, _ := hex.DecodeString(pubKey)
	copy(pub[:], pubkey)

	var myKP crypto.Keypair
	if err := myKP.Generate(); err != nil {
		fmt.Printf("Failed to generate Keypair: %s", err)
	}
	return myKP.Encrypter(pub)
}

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ignoreKeyRegex, _ := regexp.Compile("^_.*")
		encryptedRegex, _ := regexp.Compile("^EJ\\[.*\\]")
		decryptedRegex, _ := regexp.Compile("(?i)decrypted")

		file := args[0]

		cfg, err := ini.Load(file)
		if err != nil {
			fmt.Printf("Fail to read file %s: %v", args[0], err)
			return
		}

		pubKey, err := cfg.Section("").GetKey("_public_key")
		if err != nil {
			fmt.Printf("Couldn't read public key from ini: %s", err)
			return
		}

		encrypter := prepareEncrypter(pubKey.Value())

		for _, sec := range cfg.SectionStrings() {
			section, err := cfg.GetSection(sec)
			if err != nil {
				fmt.Printf("Failed parsing ini section %s\n", sec)
				return
			}
			for _, key := range section.KeyStrings() {
				if !ignoreKeyRegex.MatchString(key) {
					val := section.Key(key)
					if !decryptedRegex.MatchString(val.Comment) {
						v := val.Value()
						if !encryptedRegex.MatchString(v) {
							encrypted, err := encrypter.Encrypt([]byte(v))
							if err != nil {
								fmt.Printf("Failed encrypting key: %s\n", key)
								return
							}
							val.SetValue(fmt.Sprintf("%s", encrypted))
						}
					}
				}
			}
		}

		var output io.Writer
		if write {
			f, err := os.Create(file)
			if err != nil {
				fmt.Printf("Unable to open file %s for writing: %s\n", file, err)
				return
			}
			output = f
			defer f.Close()
		} else {
			output = os.Stdout
		}

		_, err = cfg.WriteTo(output)

		if err != nil {
			fmt.Printf("Failed to write config: %s\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVarP(&write, "write", "w", false, "write back to original file")
}
