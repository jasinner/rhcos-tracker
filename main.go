package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/containers/buildah"
	"github.com/containers/buildah/pkg/unshare"
	"github.com/containers/image/types"
	"github.com/containers/storage"
	"k8s.io/klog"
)

func getRPMs(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("tag")
	//Could improve this with skopeo inspect to get latest tag
	if tag == "" {
		http.Error(w, "Please supply a tag parameter", http.StatusBadRequest)
		return
	}
	image := fmt.Sprintf("quay.io/openshift-release-dev/ocp-release:%s", tag)
	fmt.Fprintf(w, "Getting RPMs for image %s", image)

	buildStoreOptions, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	buildStore, err := storage.GetStore(buildStoreOptions)

	builderOpts := buildah.BuilderOptions{
		FromImage:        image,                   // Starting image
		Isolation:        buildah.IsolationChroot, // Isolation environment
		CommonBuildOpts:  &buildah.CommonBuildOptions{},
		ConfigureNetwork: buildah.NetworkDefault,
		SystemContext:    &types.SystemContext{},
	}

	// getContext returns a context.TODO
	builder, err := buildah.NewBuilder(getContext(), buildStore, builderOpts)
	if err != nil {
		log.Println(fmt.Errorf("error creating buildah builder: %v", err))
	}

	mountPath, err := builder.Mount("")
	defer func() {
		err := builder.Unmount()
		if err != nil {
			klog.Errorf("failed to unmount: %v", err)
		}
	}()
	if err != nil {
		log.Println(fmt.Errorf("error mounting image content from image %s: %v", image, err))
	}
	//var out bytes.Buffer
	io.WriteString(w, mountPath)
}

func getContext() context.Context {
	return context.TODO()
}

func handleRequests() {
	http.HandleFunc("/", getRPMs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
