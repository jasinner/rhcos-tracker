package releases

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type parseTest = struct {
	input string
	want  []OpenShiftVersion
	err   error
}

// Get Cincinnati data from file
func get_mock_page(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	jsonFile, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue, nil
}

func (version *OpenShiftVersion) Equals(other *OpenShiftVersion) bool {
	return reflect.DeepEqual(version, other)
}

func findVersion(slice []OpenShiftVersion, val OpenShiftVersion) (int, bool) {
	for i, item := range slice {
		if item.Equals(&val) {
			return i, true
		}
	}
	return -1, false
}

var ocp482 = OpenShiftVersion{
	Version: "4.8.2",
	Image:   "quay.io/openshift-release-dev/ocp-release@sha256:0e82d17ababc79b10c10c5186920232810aeccbccf2a74c691487090a2c98ebc",
	Errata:  "https://access.redhat.com/errata/RHSA-2021:2438",
}

var ocp472 = OpenShiftVersion{
	Version: "4.7.2",
	Image:   "quay.io/openshift-release-dev/ocp-release@sha256:83fca12e93240b503f88ec192be5ff0d6dfe750f81e8b5ef71af991337d7c584",
	Errata:  "https://access.redhat.com/errata/RHBA-2021:0749",
}

func TestParseCincinatti(t *testing.T) {

	var tests = []parseTest{
		{"testdata/482.json", []OpenShiftVersion{ocp482}, nil},
		{"testdata/482_with_edges.json", []OpenShiftVersion{ocp482}, nil},
		{"testdata/482_472.json", []OpenShiftVersion{ocp482, ocp472}, nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.input)
		t.Run(testname, func(t *testing.T) {
			ans, err := ParseCincinnati(tt.input, get_mock_page)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("got unexpected error %v", err)
				return
			}
			if len(tt.want) != len(ans) {
				t.Errorf("Expect %v results, got %v", len(tt.want), len(ans))
				return
			}
			for _, expect := range tt.want {
				_, match := findVersion(ans, expect)
				if !match {
					t.Errorf("Expected %v, but not found in result", expect)
				}
			}
		})
	}

}
