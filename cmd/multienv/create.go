package main

import (
	"fmt"
	"github.com/fission/fission/crd"
	fv1 "github.com/fission/fission/pkg/apis/fission.io/v1"
	// "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cmdCreate = cli.Command{
	Name:        "create",
	Description: "create a multi zonal environment",
	Aliases:     []string{"mk"},
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name: "zone",
		},
		cli.StringFlag{
			Name: "runtime",
		},
		cli.StringFlag{
			Name: "name",
		},
	},
	Action: commandContext(func(ctx Context) error {
		// if !ctx.Args().Present() {
		// 	logrus.Fatal("Usage: fission-multienv create --name <name> --runtime <runtime> -- zones <zones>")
		// }
		client := getClient(ctx)
		zones := ctx.StringSlice("zone")
		runtime := ctx.String("runtime")
		name := ctx.String("name")
		crds := createMultizoneCRD(name, zones, runtime)
		for _, crd := range crds {
			_, err := client.EnvironmentCreate(crd)
			// abort if a creation fails TODO add rollback mechanism
			// to clean up environments created before failure
			if err != nil {
				return err
			}
		}
		fmt.Println("Multizone environment created!")
		return nil

	}),
}

func createMultizoneCRD(name string, zones []string, runtime string) []*crd.Environment {
	fmt.Printf("%+v", zones)
	fmt.Printf("%T", zones)
	fmt.Printf("Zone (1): %v", zones[0])
	crds := []*crd.Environment{}
	for i, _ := range zones {
		//For now keep poolsize of each 'zone' as 1.
		// TODO divide total pool size among all of the zones
		fmt.Printf("Zone: %s \n", zones[i])
		crd := createEnvironmentCRD(fmt.Sprintf("%s-%s", name, zones[i]), runtime, 1)
		crds = append(crds, crd)
	}
	return crds
}

func createEnvironmentCRD(name string, runtime string, poolSize int) *crd.Environment {
	return &crd.Environment{
		Metadata: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: fv1.EnvironmentSpec{
			Runtime: fv1.Runtime{
				Image: runtime,
			},
			Poolsize: poolSize,
		},
	}
}
