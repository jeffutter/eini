package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"bufio"
	"encoding/hex"
	"github.com/Shopify/ejson/crypto"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
)

var output string
var privKey string

func flatten_yaml(val interface{}) map[string]string {
	var output = make(map[string]string)
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		var m = reflect.ValueOf(val).Interface().(map[interface{}]interface{})
		for k, v := range m {
			var key string
			switch reflect.TypeOf(k).Kind() {
			case reflect.String:
				key = reflect.ValueOf(k).Interface().(string)
			}

			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				var s = reflect.ValueOf(v).Interface().(string)
				output[key] = s
			case reflect.Map:
				o := flatten_yaml(reflect.ValueOf(v).Interface())
				for sk, sv := range o {
					output[key+"_"+sk] = sv
				}
			default:
				fmt.Printf("Unknown Type: %v", reflect.TypeOf(v))
			}
		}
	}
	return output
}

func out(data map[string]string) {
	var priv [32]byte
	var pub [32]byte

	encryptedRegex, _ := regexp.Compile("^EJ\\[.*\\]")
	ignoreKeyRegex, _ := regexp.Compile("^_.*")

	pubKey := data["_public_key"]

	pubkey, _ := hex.DecodeString(pubKey)
	privkey, _ := hex.DecodeString(privKey)

	copy(pub[:], pubkey)
	copy(priv[:], privkey)

	myKP := crypto.Keypair{
		Public:  pub,
		Private: priv,
	}
	decrypter := myKP.Decrypter()

	for k, v := range data {
		if !ignoreKeyRegex.MatchString(k) {
			if encryptedRegex.MatchString(v) {
				decrypted, err := decrypter.Decrypt([]byte(v))
				if err != nil {
					fmt.Printf("Decryption Error: %v - %v\n", err, v)
				}
				fmt.Printf("declare -x \"%s\"=\"%s\"\n", k, decrypted)
			} else {
				fmt.Printf("declare -x \"%s\"=\"%s\"\n", k, v)
			}
		}
	}
}

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
		reader := bufio.NewReader(os.Stdin)
		privKey, _ = reader.ReadString('\n')

		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Print(err)
		}

		var val interface{}
		err = yaml.Unmarshal(b, &val)
		if err != nil {
			fmt.Println("unmarshal []byte to yaml failed: " + err.Error())
		}
		out(flatten_yaml(val))
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&output, "output", "o", "env", "output format: [env, yaml]")
}
