package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/lib-wires/wirescli"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "wires",
	Short:         "upload a file",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run:           runRoot,
}

var rootFlags struct {
	ip   string
	port string
	user string
	pass string
	file string
}

func runRoot(cmd *cobra.Command, args []string) {
	p, _ := strconv.Atoi(rootFlags.port)
	cli := wirescli.New(rootFlags.ip, p)
	fmt.Println("ip:", rootFlags.ip, "port:", p)
	path, err := os.Getwd()
	fmt.Println("current path", path)
	fmt.Println("current path", filepath.FromSlash(path))
	fmt.Println("try and open file", filepath.FromSlash(rootFlags.file))
	file, err := ioutil.ReadFile(filepath.FromSlash(rootFlags.file))
	if err != nil {
		fmt.Println(err)
		return
	}
	body := &wirescli.NodesBody{}
	var nodes interface{}
	if err := json.Unmarshal(file, &nodes); err != nil {
		fmt.Println(err)
		return
	}
	body.Nodes = nodes
	body.Pos = []float64{-1250, -1600}
	token, _ := cli.GetToken(&wirescli.WiresTokenBody{Username: rootFlags.user, Password: rootFlags.pass})
	if token == nil {
		fmt.Println("token body is nil")
		return
	}
	if token.Token == "" {
		fmt.Println("token is nil")
		return
	}
	fmt.Println("token", token.Token)
	ok, _ := cli.Upload(body, token.Token)
	if ok {
		fmt.Println("uploaded ok")
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//color.Magenta(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pFlagSet := rootCmd.PersistentFlags()
	pFlagSet.StringVarP(&rootFlags.ip, "ip", "", "localhost", "ip")
	pFlagSet.StringVarP(&rootFlags.port, "port", "", "1313", "port")
	pFlagSet.StringVarP(&rootFlags.user, "user", "", "admin", "username")
	pFlagSet.StringVarP(&rootFlags.pass, "pass", "", "admin", "password")
	pFlagSet.StringVarP(&rootFlags.file, "file", "", "../backup.json", "backup file")
}
