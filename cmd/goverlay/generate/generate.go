package generate

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/sivchari/goverlay/internal"
)

var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate is a tool to generate overlay.json",
	Long:  "generate is a tool to generate overlay.json",
	RunE:  runGenerate,
}

var (
	dist string
)

func init() {
	Cmd.Flags().StringVarP(&dist, "dist", "d", "overlay.json", "Path to the overlay.json file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	config := cmd.Flag("config").Value.String()
	b, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	var cfg internal.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}
	replaces := make([]string, len(cfg.Layers))
	for i, layer := range cfg.Layers {
		from := layer.From
		dist := layer.Dist
		replaces[i] = fmt.Sprintf(`"%s": "%s"`, from, dist)
	}
	r := strings.Join(replaces, ",\n    ")
	t := fmt.Sprintf(tmpl, r)
	var w io.Writer
	if testing.Testing() {
		w = os.Stdout
	} else {
		f, err := os.Create(dist)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	if _, err := w.Write([]byte(t)); err != nil {
		return err
	}
	return nil
}

type Replaces struct {
	Replaces []string `yaml:"Replaces"`
}

var tmpl = `{
  "Replaces": {
    %s
  }
}
`
