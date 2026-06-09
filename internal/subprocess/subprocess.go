package subprocess

import (
	"os"
	"os/exec"
)

func LaunchChildProcessWithEnvVars(args []string, proxyAddr, certPath string) (*exec.Cmd, error) {

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = append(
		os.Environ(),
		"HTTP_PROXY=http://"+proxyAddr,
		"HTTPS_PROXY=http://"+proxyAddr,
		"REQUESTS_CA_BUNDLE="+certPath,
		"SSL_CERT_FILE="+certPath,
		"NODE_EXTRA_CA_CERTS="+certPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd, cmd.Start()
}
