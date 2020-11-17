// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// evm executes EVM code snippets.
package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"

	// stage1-substate: import evm/research
	"github.com/ethereum/go-ethereum/cmd/evm/research"
)

var gitCommit = "" // Git SHA1 commit hash of the release (set via linker flags)

var (
	app = utils.NewApp(gitCommit, "the evm command line interface")

	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "output full trace logs",
	}
	MemProfileFlag = cli.StringFlag{
		Name:  "memprofile",
		Usage: "creates a memory profile at the given path",
	}
	CPUProfileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "creates a CPU profile at the given path",
	}
	StatDumpFlag = cli.BoolFlag{
		Name:  "statdump",
		Usage: "displays stack and heap memory information",
	}
	CodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "EVM code",
	}
	CodeFileFlag = cli.StringFlag{
		Name:  "codefile",
		Usage: "file containing EVM code",
	}
	GasFlag = cli.Uint64Flag{
		Name:  "gas",
		Usage: "gas limit for the evm",
		Value: 10000000000,
	}
	PriceFlag = utils.BigFlag{
		Name:  "price",
		Usage: "price set for the evm",
		Value: new(big.Int),
	}
	ValueFlag = utils.BigFlag{
		Name:  "value",
		Usage: "value set for the evm",
		Value: new(big.Int),
	}
	DumpFlag = cli.BoolFlag{
		Name:  "dump",
		Usage: "dumps the state after the run",
	}
	InputFlag = cli.StringFlag{
		Name:  "input",
		Usage: "input for the EVM",
	}
	VerbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "sets the verbosity level",
	}
	CreateFlag = cli.BoolFlag{
		Name:  "create",
		Usage: "indicates the action should be create rather than call",
	}
	DisableGasMeteringFlag = cli.BoolFlag{
		Name:  "nogasmetering",
		Usage: "disable gas metering",
	}
	GenesisFlag = cli.StringFlag{
		Name:  "prestate",
		Usage: "JSON file with prestate (genesis) config",
	}
	MachineFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output trace logs in machine readable format (json)",
	}
	SenderFlag = cli.StringFlag{
		Name:  "sender",
		Usage: "The transaction origin",
	}
	DisableMemoryFlag = cli.BoolFlag{
		Name:  "nomemory",
		Usage: "disable memory output",
	}
	DisableStackFlag = cli.BoolFlag{
		Name:  "nostack",
		Usage: "disable stack output",
	}
)

// stage1-substate: t8n-substate command
var stateTransitionSubstateCommand = cli.Command{
	Action:    research.TransitionSubstate,
	Name:      "t8n-substate",
	Aliases:   []string{"t8n-substate"},
	Usage:     "executes full state transitions and check output consistency",
	ArgsUsage: "<blockNumFirst> <blockNumLast>",
	Flags: []cli.Flag{
		research.WorkersFlag,
		research.SkipTransferTxsFlag,
		research.SkipCallTxsFlag,
		research.SkipCreateTxsFlag,
	},
	Description: `
 The transition-substate (t8n-substate) command requires
 two arguments: <blockNumFirst> <blockNumLast>
 <blockNumFirst> and <blockNumLast> are the first and
 last block of the inclusive range of blocks to replay transactions.`,
}

// stage1-substate: dump-substate command
var dumpSubstateCommand = cli.Command{
	Action:    research.DumpSubstate,
	Name:      "dump-substate",
	Usage:     "dump a range of substates into target LevelDB",
	ArgsUsage: "<targetPath> <blockNumFirst> <blockNumLast>",
	Flags: []cli.Flag{
		research.WorkersFlag,
	},
	Description: `
 The dump-substate command requires three arguments:
 <targetPath> <blockNumFirst> <blockNumLast>
 <targetPath> is the target LevelDB where to dump substate.
 <blockNumFirst> and <blockNumLast> are the first and
 last block of the inclusive range of blocks to replay transactions.`,
}

// stage1-substate: size-substate command
var sizeSubstateCommand = cli.Command{
	Action:    research.SizeSubstate,
	Name:      "size-substate",
	Usage:     "calculate size of decompressed values in substate DB",
	ArgsUsage: "<blockNumFirst> <blockNumLast>",
	Flags:     []cli.Flag{},
	Description: `
 The size-substate command requires two arguments:
 <blockNumFirst> <blockNumLast>
 <blockNumFirst> and <blockNumLast> are the first and
 last block of the inclusive range of blocks to replay transactions.`,
}

// stage1-substate: replay-fork command
var replayForkCommand = cli.Command{
	Action:    research.ReplayFork,
	Name:      "replay-fork",
	Usage:     "executes and check output consistency of all transactions in the range with the given hard-fork",
	ArgsUsage: "<blockNumFirst> <blockNumLast>",
	Flags: []cli.Flag{
		research.WorkersFlag,
		research.SkipTransferTxsFlag,
		research.SkipCallTxsFlag,
		research.SkipCreateTxsFlag,
		research.HardForkFlag,
	},
	Description: `
 The replay-fork command requires two arguments:
 <blockNumFirst> <blockNumLast>

 <blockNumFirst> and <blockNumLast> are the first and
 last block of the inclusive range of blocks to replay transactions.
 --hard-fork parameter is recommended for this command.`,
}

func init() {
	app.Flags = []cli.Flag{
		CreateFlag,
		DebugFlag,
		VerbosityFlag,
		CodeFlag,
		CodeFileFlag,
		GasFlag,
		PriceFlag,
		ValueFlag,
		DumpFlag,
		InputFlag,
		DisableGasMeteringFlag,
		MemProfileFlag,
		CPUProfileFlag,
		StatDumpFlag,
		GenesisFlag,
		MachineFlag,
		SenderFlag,
		DisableMemoryFlag,
		DisableStackFlag,
	}
	app.Commands = []cli.Command{
		compileCommand,
		disasmCommand,
		runCommand,
		// stage1-substate: transition-substate (t8n-substate) command
		stateTransitionSubstateCommand,
		// stage1-substate: dump-substate command
		dumpSubstateCommand,
		// stage1-substate: size-substate command
		sizeSubstateCommand,
		// stage1-substate: replay-fork command
		replayForkCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
