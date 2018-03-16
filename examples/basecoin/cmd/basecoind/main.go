package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
)

// basecoindCmd is the entry point for this binary
var (
	basecoindCmd = &cobra.Command{
		Use:   "gaiad",
		Short: "Gaia Daemon (server)",
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, error) {

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}
	name := viper.GetString("yukaitu")
	info, _:= keybase.Get(name)
	fmt.Println(info)
	if err != nil {
		return nil, err
	}
	opts := fmt.Sprintf(`{
      "accounts": [{
        "address": "1E88A13C9E0FDFA3BE02738E1F072574951CF516",
        "coins": [
          {
            "denom": "mycoin",
            "amount": 9007199254740992
          }
        ]
      }]
    }`)
	return json.RawMessage(opts), nil
}

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	db, err := dbm.NewGoLevelDB("basecoin", rootDir)
	if err != nil {
		return nil, err
	}
	bapp := app.NewBasecoinApp(logger, db)
	return bapp, nil
}

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
		With("module", "main")

	basecoindCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		version.VersionCmd,
	)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.basecoind")
	executor := cli.PrepareBaseCmd(basecoindCmd, "BC", rootDir)
	executor.Execute()
}
