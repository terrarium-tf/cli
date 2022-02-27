package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ojizero/gofindup"
	"github.com/spf13/cobra"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func ArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("requires a workspace and a stack path")
	}
	if _, err := os.Stat(args[1]); os.IsNotExist(err) {
		return fmt.Errorf("invalid path given: %s", args[1])
	}
	return nil
}

func Vars(cmd cobra.Command, env string, stackPath string) ([]string, map[string]interface{}) {
	// collect global vars
	vars := make(map[string]interface{})
	var files []string

	// collect global vars
	f := readVarsFile(cmd, "global.tfvars.json", stackPath, &vars)
	if f != "" {
		files = append(files, f)
	}

	// collect global env vars
	f = readVarsFile(cmd, fmt.Sprintf("%s.tfvars.json", strings.ToLower(env)), stackPath+"/../", &vars)
	if f != "" {
		files = append(files, f)
	}

	// collect stack global vars
	f = readVarsFile(cmd, "app.tfvars.json", stackPath, &vars)
	if f != "" {
		files = append(files, f)
	}

	// collect stack env vars
	f = readVarsFile(cmd, fmt.Sprintf("%s.tfvars.json", strings.ToLower(env)), stackPath, &vars)
	if f != "" {
		files = append(files, f)
	}

	verbose, _ := cmd.Parent().PersistentFlags().GetBool("verbose")
	if verbose {
		if len(files) > 0 {
			cmd.Printf(InfoColorLine, "Collected vars files:")
			for _, f := range files {
				cmd.Printf(WarningColorLine, f)
			}
		}

		cmd.Println("")
		cmd.Printf(InfoColorLine, "Collected vars:")
		maxlen := 0
		for k, _ := range vars {
			if len(k) > maxlen {
				maxlen = len(k)
			}
		}

		for k, v := range vars {
			cmd.Printf(WarningColorMap, maxlen, k, VarToString(v))
		}
	}
	return files, vars
}

func VarToString(v interface{}) string {
	strVal := ""
	switch t := v.(type) {
	case int:
		strVal = fmt.Sprintf("%d", t)
	case float64:
		_int, _float := math.Modf(t)
		if _float == 0 {
			strVal = fmt.Sprintf("%d", int(_int))
		} else {
			strVal = fmt.Sprintf("%f", t)
		}
	case bool:
		strVal = fmt.Sprintf("%t", t)
	default:
		strVal = fmt.Sprintf("%s", t)
	}
	return strVal
}

func readVarsFile(cmd cobra.Command, name string, path string, vars *map[string]interface{}) string {
	file, err := gofindup.FindupFrom(name, path)
	if err != nil {
		log.Fatal(err)
	}

	if file != "" {
		content, err := os.ReadFile(file)
		if err != nil {
			cmd.PrintErr("error reading file ", file, err)
			os.Exit(1)
		}
		err = json.Unmarshal(content, vars)
		if err != nil {
			cmd.PrintErr("error reading json file ", file, err)
			os.Exit(1)
		}

		absPath, err := filepath.Abs(file)
		if err != nil {
			cmd.PrintErr("error reading file ", file, err)
			os.Exit(1)
		}
		return absPath
	}

	return ""
}

func GetVar(name string, cmd cobra.Command, mergedVars map[string]interface{}, required bool) string {
	var _var string
	flag := cmd.Flags().Lookup(fmt.Sprintf("state-%s", name))

	if flag != nil && flag.Changed {
		_var = flag.Value.String()
	} else {
		if fileVar, ok := mergedVars[name]; ok {
			_var = VarToString(fileVar)
		}
	}

	if required && _var == "" {
		cmd.PrintErrf(ErrorColorLine, fmt.Sprintf("unable to configure remote state, '%s' was not found in var files and not provided with '-state-%s'", name, name))
		os.Exit(1)
	}

	return _var
}
