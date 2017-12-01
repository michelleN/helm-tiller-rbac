package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
)

const globalUsage = `
Enable an RBAC profile for Tiller
An RBAC profile at the moment is a Helm chart that contains a set of Kubernetes manifests defining roles, role-bindings, service accounts that can be templated

Example Usage: 
   $ helm secure-tiller dev-team-rbac-profile-chart/
`

var namespace string
var version = "DEV"

func main() {
	cmd := &cobra.Command{
		Use:   "secure-tiller [RBAC_PROFILE]",
		Short: globalUsage,
		RunE:  run,
	}

	f := cmd.Flags()
	f.StringVar(&namespace, "namespace", "default", "namespace of Tiller to apply profile")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing argument: path to chart containing rbac profile")
	}

	if len(args) > 1 {
		return errors.New("TMI")
	}
	profilePath := args[0]
	profileName := filepath.Base(profilePath)

	var outb bytes.Buffer
	var errb bytes.Buffer

	_, err := chartutil.Load(profilePath)
	if err != nil {
		return err
	}

	renderTemplates := exec.Command("helm", "template", "--namespace", namespace, profilePath)
	renderTemplates.Stdout = &outb
	renderTemplates.Stderr = &errb
	err = renderTemplates.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	manifests := []byte(outb.String())

	dir, err := ioutil.TempDir("", "secure-tiller")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	manifestsFile := filepath.Join(dir, "profile")
	if err := ioutil.WriteFile(manifestsFile, manifests, 0666); err != nil {
		return err
	}

	outb.Reset()
	errb.Reset()
	createSvcAccount := exec.Command("kubectl", "create", "serviceaccount", profileName, "--namespace", namespace)
	createSvcAccount.Stdout = &outb
	createSvcAccount.Stderr = &errb
	err = createSvcAccount.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	fmt.Print(outb.String())

	outb.Reset()
	errb.Reset()
	createRoleBindingBits := exec.Command("kubectl", "apply", "-f", manifestsFile)
	createRoleBindingBits.Stdout = &outb
	createRoleBindingBits.Stderr = &errb
	err = createRoleBindingBits.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	fmt.Println(outb.String())
	outb.Reset()
	errb.Reset()

	specPatch := `{"spec":{"template":{"spec":{"serviceAccount":"` + profileName + `","serviceAccountName":"` + profileName + `"}}}}`

	outb.Reset()
	errb.Reset()
	// this command will terminate & retstart pods automatically after attaching the service account
	attachSvcAccount := exec.Command("kubectl", "patch", "--namespace", namespace, "deployment", "tiller-deploy", "-p", specPatch)
	attachSvcAccount.Stdout = &outb
	attachSvcAccount.Stderr = &errb

	err = attachSvcAccount.Run()
	if err != nil {
		return errors.New(errb.String())
	}

	fmt.Println("Congrats! Your Tiller in the " + namespace + " namespace is secured with the " + profileName + " service account\nYou can verify your Tiller config by running this command\n\t$ kubectl -n " + namespace + " get deployment tiller-deploy -o json\n")

	//TODO: verify that the service account you wanted is now part of the tiller deployment spec

	return nil

}
