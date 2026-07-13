// Portions of the trust-store installation logic in this file
// were adapted from mkcert:
// https://github.com/FiloSottile/mkcert
//
// mkcert is licensed under the BSD 3-Clause License.

package truststore

import (
	"github.com/hrpofficial736/promtrace/internal/logger"
	"os/exec"
)

type TrustStoreInstaller interface {
	Install(certPath string) error
	Uninstall(certPath string) error
}

// truststore management for macOS
type DarwinInstaller struct{}

func (d DarwinInstaller) Install(certPath string) error {
	cmd := exec.Command(
		"sudo",
		"security",
		"add-trusted-cert",
		"-d",
		"-k",
		"/Library/Keychains/System.keychain",
		certPath,
	)

	return cmd.Run()
}

func (d DarwinInstaller) Uninstall(certPath string) error {
	cmd := exec.Command(
		"sudo",
		"security",
		"remove-trusted-cert",
		"-d",
		certPath,
	)

	return cmd.Run()
}

// truststore management for linux

type LinuxInstaller struct{}

func (l LinuxInstaller) Install(certPath string) error {
	// Step 1: ensure the directory exists
	cmd1 := exec.Command("sudo", "mkdir", "-p", "/usr/local/share/ca-certificates")
	if err := cmd1.Run(); err != nil {
		logger.Log.Error("error while executing mkdir command for installation", "error", err)
		return err
	}

	// Step 2: copy the cert
	cmd2 := exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/promtrace-ca.crt")
	if err := cmd2.Run(); err != nil {
		logger.Log.Error("error while executing cp command for installation", "error", err)
		return err
	}

	// Step 3: update-ca-certificates
	cmd3 := exec.Command("sudo", "update-ca-trust")
	return cmd3.Run()
}

func (l LinuxInstaller) Uninstall(certPath string) error {
	cmd := exec.Command(
		"sudo",
		"rm",
		"/usr/local/share/ca-certificates/ca.crt",
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	return exec.Command("sudo", "update-ca-certificates").Run()
}

// truststore management for windows

type WindowsInstaller struct{}

func (w WindowsInstaller) Install(certPath string) error {
	cmd := exec.Command(
		"certutil",
		"-addstore",
		"Root",
		certPath,
	)

	return cmd.Run()
}

func (w WindowsInstaller) Uninstall(certPath string) error {
	cmd := exec.Command(
		"certutil",
		"-delstore",
		"Root",
		certPath,
	)

	return cmd.Run()
}
