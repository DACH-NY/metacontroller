/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/0xRLG/ocworkqueue"
	"github.com/golang/glog"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"metacontroller.app/apis/metacontroller/v1alpha1"
	mcclientset "metacontroller.app/client/generated/clientset/internalclientset"
	mcinformers "metacontroller.app/client/generated/informer/externalversions"
	"metacontroller.app/controller/composite"
	"metacontroller.app/controller/decorator"
	dynamicclientset "metacontroller.app/dynamic/clientset"
	dynamicdiscovery "metacontroller.app/dynamic/discovery"
	dynamicinformer "metacontroller.app/dynamic/informer"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var (
	discoveryInterval = flag.Duration("discovery-interval", 30*time.Second, "How often to refresh discovery cache to pick up newly-installed resources")
	informerRelist    = flag.Duration("cache-flush-interval", 30*time.Minute, "How often to flush local caches and relist objects from the API server")
	debugAddr         = flag.String("debug-addr", ":9999", "The address to bind the debug http endpoints")
	clientConfigPath  = flag.String("client-config-path", "", "Path to kubeconfig file (same format as used by kubectl); if not specified, use in-cluster config")
	namespace = flag.String("namespace", "default", "Namespace where the controller listens")
	)

type controller interface {
	Start()
	Stop()
}

func main() {
	flag.Parse()

	glog.Infof("Discovery cache flush interval: %v", *discoveryInterval)
	glog.Infof("API server object cache flush interval: %v", *informerRelist)
	glog.Infof("Debug http server address: %v", *debugAddr)

	var config *rest.Config
	var err error
	if *clientConfigPath != "" {
		glog.Infof("Using current context from kubeconfig file: %v", *clientConfigPath)
		config, err = clientcmd.BuildConfigFromFlags("", *clientConfigPath)
	} else {
		glog.Info("No kubeconfig file specified; trying in-cluster auto-config...")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatal(err)
	}

	// Periodically refresh discovery to pick up newly-installed resources.
	dc := discovery.NewDiscoveryClientForConfigOrDie(config)
	resources := dynamicdiscovery.NewResourceMap(dc)
	// We don't care about stopping this cleanly since it has no external effects.
	resources.Start(*discoveryInterval)

	// Create informerfactory for metacontroller api objects.
	mcClient, err := mcclientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Can't create client for api %s: %v", v1alpha1.SchemeGroupVersion, err)
	}
	// namespace must be supplied via an external parameter
	// we've hardcoded the names space as "todo"
	nsOpt := mcinformers.WithNamespace(*namespace)
	mcInformerFactory := mcinformers.NewSharedInformerFactoryWithOptions(mcClient, *informerRelist, nsOpt)

	// Create dynamic clientset (factory for dynamic clients).
	dynClient, err := dynamicclientset.New(config, resources)
	if err != nil {
		glog.Fatal(err)
	}

	// Create dynamic informer factory (for sharing dynamic informers).
	dynInformers := dynamicinformer.NewSharedInformerFactory(dynClient, *informerRelist, *namespace)

	workqueue.SetProvider(ocworkqueue.MetricsProvider())
	view.Register(ocworkqueue.DefaultViews...)

	// Start metacontrollers (controllers that spawn controllers).
	// Each one requests the informers it needs from the factory.
	controllers := []controller{
		composite.NewMetacontroller(resources, dynClient, dynInformers, mcInformerFactory, mcClient),
		decorator.NewMetacontroller(resources, dynClient, dynInformers, mcInformerFactory),
	}

	// Start all requested informers.
	// We don't care about stopping this cleanly since it has no external effects.
	mcInformerFactory.Start(nil)

	// Start all controllers.
	for _, c := range controllers {
		c.Start()
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		glog.Fatalf("can't create prometheus exporter: %v", err)
	}
	view.RegisterExporter(exporter)

	mux := http.NewServeMux()
	mux.Handle("/metrics", exporter)
	srv := &http.Server{
		Addr:    *debugAddr,
		Handler: mux,
	}
	go func() {
		glog.Errorf("Error serving debug endpoint: %v", srv.ListenAndServe())
	}()

	// On SIGTERM, stop all controllers gracefully.
	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigchan
	glog.Infof("Received %q signal. Shutting down...", sig)

	var wg sync.WaitGroup
	for _, c := range controllers {
		wg.Add(1)
		go func(c controller) {
			defer wg.Done()
			c.Stop()
		}(c)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Shutdown(context.Background())
	}()
	wg.Wait()
}
