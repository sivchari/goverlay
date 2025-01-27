package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

func goverlayCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "goverlay is a tool to generate overlayed code",
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			return func(state *script.State) (string, string, error) {
				r, w, err := os.Pipe()
				if err != nil {
					return "", "", err
				}

				// Redirect stdout to a pipe.
				stdout := os.Stdout
				os.Stdout = w

				cmdArgs := append([]string{"goverlay"}, args...)
				os.Args = cmdArgs
				main()

				// restore stdout
				os.Stdout = stdout
				w.Close()

				// Read the output from the pipe.
				var buf bytes.Buffer
				io.Copy(&buf, r)

				s.Logf("args: %s, stdout: %s\n", os.Args, buf.String())

				return buf.String(), "", nil
			}, nil
		},
	)
}

func generateCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "generate is a tool to generate overlay.json",
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			return func(state *script.State) (string, string, error) {
				r, w, err := os.Pipe()
				if err != nil {
					return "", "", err
				}

				// Redirect stdout to a pipe.
				stdout := os.Stdout
				os.Stdout = w

				cmd := []string{"goverlay", "generate"}
				cmdArgs := append(cmd, args...)
				os.Args = cmdArgs
				main()

				// restore stdout
				os.Stdout = stdout
				w.Close()

				// Read the output from the pipe.
				var buf bytes.Buffer
				io.Copy(&buf, r)

				s.Logf("args: %s, stdout: %s\n", os.Args, buf.String())

				return buf.String(), "", nil
			}, nil
		},
	)
}

func TestScript(t *testing.T) {
	ctx := context.Background()
	engine := script.NewEngine()
	engine.Cmds["goverlay"] = goverlayCmd()
	engine.Cmds["generate"] = generateCmd()
	env := os.Environ()
	scripttest.Test(t, ctx, engine, env, "testdata/*.txt")
}
