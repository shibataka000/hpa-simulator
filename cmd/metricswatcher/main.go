package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/shibataka000/metrics-watcher/pkg/kubernetes"
	"github.com/shibataka000/metrics-watcher/pkg/metricswatcher"
	"github.com/urfave/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func action(c *cli.Context) error {
	config, err := kubernetes.NewClientConfig(c.String("kubeconfig"), c.String("context"))
	if err != nil {
		return err
	}

	podQuery := c.Args().Get(0)
	podQueryRegex, err := regexp.Compile(podQuery)
	if err != nil {
		return err
	}

	metricswatcher.Watch(config, c.String("namespace"), podQueryRegex)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "metricswatcher"
	app.Usage = "Output some information to know HorizontalPodAutoscaler internal behavior"
	app.UsageText = "metricswatcher [pod_query]"
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
