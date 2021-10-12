package main

import (
	"bytes"
	"fmt"
	"io"
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
	imageTag := fmt.Sprintf("quay.io/openshift-release-dev/ocp-release:%s", tag)
	fmt.Fprintf(w, "Getting RPMs for image %s", imageTag)
	cmd := exec.Command("/usr/bin/oc", "adm", "release", "info", imageTag)
	var stderr bytes.Buffer
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	io.WriteString(w, out.String())
}

func handleRequests() {
	http.HandleFunc("/", getRPMs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
