package main

import (
	"context"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
	"github.com/syncromatics/go-kit/v2/log"
	"github.com/thzinc/mdns-ingress-publisher/internal/controller"
	"github.com/ugjka/mdns"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mdns-ingress-publisher",
		Short: "Watches for Kubernetes Ingress object events and publishes mDNS records accordingly",
		RunE: func(*cobra.Command, []string) error {
			kubeconfigPath := viper.GetString("kubeconfig")
			_, err := os.Stat(kubeconfigPath)
			if err != nil {
				return errors.Wrap(err, "failed to find Kubernetes client config")
			}

			config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				return errors.Wrap(err, "failed to interpret Kubernetes client config")
			}

			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return errors.Wrap(err, "failed to create Kubernetes client")
			}

			zone, err := mdns.New(true, false)
			if err != nil {
				return errors.Wrap(err, "failed to initialize mDNS zone")
			}
			defer zone.Shutdown()

			group := cmd.NewProcessGroup(context.Background())
			ctx := group.Context()

			ingresses := client.NetworkingV1().Ingresses("")

			defaultTTS := viper.GetInt("default-tts")
			group.Go(controller.NewWatcher(ctx, ingresses, zone, defaultTTS))

			return group.Wait()
		},
	}
)

func init() {
	var defaultKubeConfig string
	currentUser, err := user.Current()
	if err != nil {
		defaultKubeConfig = ""
	} else {
		defaultKubeConfig = filepath.Join(currentUser.HomeDir, ".kube", "config")
	}
	rootCmd.Flags().String("kubeconfig", defaultKubeConfig, "Path to the Kubernetes client config file")
	rootCmd.Flags().Int("default-tts", 60, "Default mDNS record time to live (TTL); may be individually overridden using the \"mdns-ingress-publisher/tts\" annotation")

	viper.SetEnvPrefix("AVAHI_CONTROLLER")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.BindPFlags(rootCmd.Flags())
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("failed to terminate cleanly",
			"err", err)
	}
}
