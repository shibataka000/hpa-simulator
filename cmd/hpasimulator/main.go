package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/shibataka000/hpa-simulator/pkg/hpasimulator"
	"github.com/shibataka000/hpa-simulator/pkg/kubernetes"
	"github.com/urfave/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func action(c *cli.Context) error {
	clientConfig, err := kubernetes.NewClientConfig(c.String("kubeconfig"), c.String("context"))
	if err != nil {
		return err
	}
	log.Printf("%v\n", clientConfig)

	config, err := hpasimulator.NewConfig(c.String("namespace"), c.String("selector"))
	if err != nil {
		return err
	}
	simulator, err := hpasimulator.NewHpaSimulator(clientConfig, config)
	if err != nil {
		return err
	}

	err = simulator.Start()

	return err
}

func main() {
	app := cli.NewApp()
	app.Name = "hpasimulator"
	app.Usage = "Simulate HorizontalPodAutoscaler and output some information to know it's internal behavior"
	app.UsageText = "hpasimulator [flags]"
	app.Version = "v0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "kubeconfig",
			Value: filepath.Join(os.Getenv("HOME"), ".kube", "config"),
			Usage: "Path to the kubeconfig file to use for CLI requests.",
		},
		cli.StringFlag{
			Name:  "context",
			Value: "",
			Usage: "The name of the kubeconfig context to use",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Value: "default",
			Usage: "If present, the namespace scope for this CLI request",
		},
		cli.StringFlag{
			Name:  "selector, l",
			Value: "",
			Usage: "Selector (label query) to filter on, supports '='.(e.g. -l key1=value1,key2=value2)",
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
