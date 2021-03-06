package cmd

import (
	"errors"
	"fmt"
	"github.com/jeffutter/eini/ini"
	"github.com/spf13/cobra"
	"io"
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
		checkError(err)

		pubkey, err := cfg.PubKey()
		checkError(err)

		privkey, err := getPrivateKey(os.Stdin, pubkey)
		checkError(err)

		output, err := cfg.Decrypt(pubkey, privkey)
		checkError(err)

		for _, line := range output {
			fmt.Fprintln(os.Stdout, line)
		}

	},
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func checkErrorf(err error, str string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, str, args...)
		os.Exit(1)
	}
}

func getPrivateKey(reader io.Reader, pubkey string) (string, error) {
	var privKey string

	stdinContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return privKey, fmt.Errorf("Failed to read from stdin: %s", err)
	}

	privKey = strings.TrimSpace(string(stdinContent))

	if privKey == "" {
		privKey, err = readPrivateKeyFromDisk(pubkey, keydir)
		if err != nil {
			return privKey, fmt.Errorf("Error reading private key from disk: %s", err)
		}

		if privKey == "" {
			return privKey, errors.New("Private key not found")
		}
	}
	return privKey, nil
}

func readPrivateKeyFromDisk(pubkey string, keydir string) (string, error) {
	var privkey string
	var fileContents []byte
	var err error

	keyFile := fmt.Sprintf("%s/%s", keydir, pubkey)
	fileContents, err = ioutil.ReadFile(keyFile)
	if err != nil {
		err = fmt.Errorf("Couldn't read key file (%s)", err.Error())
		return privkey, err
	}

	privkey = string(fileContents)
	return privkey, nil
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&output, "output", "o", "env", "output format: [env]")
	decryptCmd.Flags().StringVarP(&keydir, "keydir", "k", "/opt/ejson/keys", "Directory containing EJSON keys")
}
