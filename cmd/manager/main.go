// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)

	addonutils "github.com/stolostron/governance-policy-spec-sync/pkg/utils"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	"github.com/stolostron/governance-policy-spec-sync/cmd/manager/tool"
	"github.com/stolostron/governance-policy-spec-sync/pkg/apis"
	"github.com/stolostron/governance-policy-spec-sync/pkg/controller"
	"github.com/stolostron/governance-policy-spec-sync/version"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost       = "0.0.0.0"
	metricsPort int32 = 8384
)
var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Operator Version: %s", version.Version))
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

func main() {
	tool.ProcessFlags()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()

	logf.SetLogger(zap.New())

	printVersion()

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}

	// Get managedconfig to talk to hub apiserver
	var managedCfg *rest.Config
	if tool.Options.ManagedConfigFilePathName == "" {
		found := false
		tool.Options.ManagedConfigFilePathName, found = os.LookupEnv("MANAGED_CONFIG")
		if found {
			log.Info("Found ENV MANAGED_CONFIG, initializing using", "tool.Options.ManagedConfigFilePathName",
				tool.Options.ManagedConfigFilePathName)
			managedCfg, err = clientcmd.BuildConfigFromFlags("", tool.Options.ManagedConfigFilePathName)
		} else {
			managedCfg, err = config.GetConfig()
			if err != nil {
				log.Error(err, "")
				os.Exit(1)
			}
		}
	}

	// Get hubconfig to talk to hub apiserver
	if tool.Options.HubConfigFilePathName == "" {
		found := false
		tool.Options.HubConfigFilePathName, found = os.LookupEnv("HUB_CONFIG")
		if found {
			log.Info("Found ENV HUB_CONFIG, initializing using", "tool.Options.HubConfigFilePathName",
				tool.Options.HubConfigFilePathName)
		}
	}

	hubCfg, err := clientcmd.BuildConfigFromFlags("", tool.Options.HubConfigFilePathName)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Set default manager options
	options := manager.Options{
		LeaderElection:         true,
		LeaderElectionID:       "policy-spec-sync-lock",
		LeaderElectionConfig:   managedCfg,
		Namespace:              namespace,
		MetricsBindAddress:     fmt.Sprintf("%s:%d", metricsHost, metricsPort),
		HealthProbeBindAddress: tool.Options.ProbeAddr,
	}

	// Add support for MultiNamespace set in WATCH_NAMESPACE (e.g ns1,ns2)
	// Note that this is not intended to be used for excluding namespaces, this is better done via a Predicate
	// Also note that you may face performance issues when using this with a high number of namespaces.
	// More Info: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/cache#MultiNamespacedCacheBuilder
	if strings.Contains(namespace, ",") {
		options.Namespace = ""
		options.NewCache = cache.MultiNamespacedCacheBuilder(strings.Split(namespace, ","))
	}

	// Create a new manager to provide shared dependencies and start components
	mgr, err := manager.New(hubCfg, options)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	cc, err := addonutils.NewConfigChecker("policy-spec-sync", tool.Options.HubConfigFilePathName)
	if err != nil {
		log.Error(err, "unable to setup a configChecker")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", cc.Check); err != nil {
		log.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr, managedCfg); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}
