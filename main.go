package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const globalUsage = `
Enable an RBAC profile for Tiller
An RBAC profile at the moment is set of Kubernetes manifests defining roles, role-bindings, service accounts

Example Usage: 
   $ helm secure-tiller dev-team-rbac-profile/
`

var namespace string
var version = "0.0.1-dev"

func main() {
	cmd := &cobra.Command{
		Use:   "secure-tiller [RBAC_PROFILE_PATH]",
		Short: globalUsage,
		RunE:  run,
	}

	f := cmd.Flags()
	f.StringVar(&namespace, "namespace", "default", "namespace to create service account")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("path to rbac profile")
	}

	if len(args) > 1 {
		return errors.New("TMI")
	}
	profilePath := args[0]
	profileName := filepath.Base(profilePath)

	var outb bytes.Buffer
	var errb bytes.Buffer

	createSvcAccount := exec.Command("kubectl", "create", "serviceaccount", profileName, "--namespace", namespace)
	createSvcAccount.Stdout = &outb
	createSvcAccount.Stderr = &errb
	err := createSvcAccount.Run()
	if err != nil {
		return errors.New(errb.String())
	}

	createRoleBindingBits := exec.Command("kubectl", "apply", "-f", profilePath)
	createRoleBindingBits.Stdout = &outb
	createRoleBindingBits.Stderr = &errb
	err = createRoleBindingBits.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	fmt.Println(outb.String())

	specPatch := `{"spec":{"template":{"spec":{"serviceAccount":"` + profileName + `","serviceAccountName":"` + profileName + `"}}}}`

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
