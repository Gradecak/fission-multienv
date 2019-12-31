package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/fission/fission/crd"
	"github.com/urfave/cli"
)

var zones = [...]string{"-de", "-nl", "-fr"}

var cmdList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Description: "List zones available for all environments",
	Action: commandContext(func(ctx Context) error {
		client := getClient(ctx)
		resp, err := client.EnvironmentList("")
		if err != nil {
			logrus.Errorf("Error fetching environment %v", err)
		}
		envs := parseEnvironments(resp)
		table(os.Stdout, []string{"ENVIRONMENT", "ZONES"}, envs)
		return nil
	}),
}

func parseEnvironments(info []crd.Environment) [][]string {
	envs := map[string][]string{} // <original_env>:<[zone variants]>
	for _, inf := range info {
		env := inf.Metadata.Name
		for _, zone := range zones {
			if strings.HasSuffix(env, zone) {
				baseEnv := strings.Replace(env, zone, "", 1)
				envs[baseEnv] = append(envs[baseEnv], zone[1:])
			}
		}
	}

	//format for displaying in table format
	formatted := [][]string{}
	for k, v := range envs {
		row := []string{k, strings.Join(v, ",")}
		formatted = append(formatted, row)
	}
	return formatted
}
