package main

import (
	"flag"
	"log"
	"time"

	"github.com/coreos/fleet/schema"
)

var (
	destroyFlag             = flag.Bool("destroy", false, "Destroy units not found in the definition")
	fleetEndpoint           = flag.String("fleetEndpoint", "", "Fleet API http endpoint: `http://host:port`")
	socksProxy              = flag.String("socksProxy", "", "address of socks proxy, e.g., 127.0.0.1:9050")
	destroyServiceBlacklist = map[string]struct{}{"deployer.service": struct{}{}, "deployer.timer": struct{}{}}
	rootURI                 = flag.String("rootURI", "", "Base uri to use when constructing service file URI. Only used if service file URI is relative.")
)

type services struct {
	Services []service `yaml:"services"`
}

type service struct {
	Name                 string `yaml:"name"`
	Version              string `yaml:"version"`
	Count                int    `yaml:"count"`
	URI                  string `yaml:"uri"`
	DesiredState         string `yaml:"desiredState"`
	SequentialDeployment bool   `yaml:"sequentialDeployment"`
}

type serviceDefinitionClient interface {
	servicesDefinition() (services, error)
	serviceFile(service service) ([]byte, error)
}

type zddInfo struct {
	unit  *schema.Unit
	count int
}

func main() {
	flag.Parse()
	if *fleetEndpoint == "" {
		log.Fatal("Fleet endpoint is required")
	}

	if *rootURI == "" {
		log.Fatal("Services definition file uri is required")
	}
	d, err := newDeployer()
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("Starting deploy run")
		if err := d.deployAll(); err != nil {
			log.Fatalf("Failed to run deploy : %v\n", err.Error())
		}
		log.Printf("Finished deploy run")
		time.Sleep(1 * time.Minute)
	}
}