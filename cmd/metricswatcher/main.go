package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/shibataka000/metrics-watcher/pkg/metricswatcher"
	"github.com/urfave/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func action(c *cli.Context) error {
	config, err := clientcmd.BuildConfigFromFlags("", c.String("kubeconfig"))
	if err != nil {
		return err
	}
	metricswatcher.Watch(config, c.String("namespace"))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "metricswatcher"
	app.Usage = "Output some information to know HorizontalPodAutoscaler internal behavior"
	app.UsageText = "metricswatcher pod_prefix"
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
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
