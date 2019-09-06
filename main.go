package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func getRPMs(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	//Could improve this with skopeo inspect to get latest tag
	if tag == "" {
		http.Error(w, "Please supply a tag parameter", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Getting RPMs for version %s", tag)
	image_tag := fmt.Sprintf("quay.io/openshift-release-dev/ocp-release:%s", tag)
	cmd := exec.Command("/usr/bin/oc", "adm", "release", "info", "--image-for=\"\"", image_tag)
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}

func handleRequests() {
	http.HandleFunc("/", getRPMs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
