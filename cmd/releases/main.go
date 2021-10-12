package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jasinner/rhcos-tracker/pkg/releases"
)

var requestURL = "https://api.openshift.com/api/upgrades_info/v1/graph?channel=stable-4."

func getEnvOrDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func main() {
	minorMax, _ := strconv.Atoi(getEnvOrDefault("MINOR_MAX_VERSION", "8"))
	minorMin, _ := strconv.Atoi(getEnvOrDefault("MINOR_MIN_VERSION", "5"))

	for i := minorMax; i >= minorMin; i-- {
		products, err := releases.ParseCincinnati(requestURL+strconv.Itoa(i), releases.GetPage)
		if err != nil {
			fmt.Printf("Failed to get release image for OCP version: 4.%v", i)
			os.Exit(1)
		}
		//read ocp pull secret from env

		for _, p := range products {
			//make sure version starts with '4.'
			fmt.Printf("Found product %v for minor %v\n", p, i)
			//get os_sha
			//oc adm release info --image-for="" 4.7.2 | grep machine-os-content | awk '{print $2}'

			//if package list for p.Version (4.7.2) doesn't already exist in DB:
			//persist package list
			//mkdir /tmp/os-release
			//oc image extract quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:0b2c764f69eb4663efb2954e74d0c235b5edcb429fd9d66f151dc666be03f63c --path /:/tmp/os-release -a ~/pull-secret.txt
			//parse /tmp/os-release/pkglist.txt
			//remove /tmp/os-release
		}
	}
}
