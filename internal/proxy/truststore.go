package proxy

import (
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
	cmd := exec.Command(
		"sudo",
		"cp",
		certPath,
		"/usr/local/share/ca-certificates/myca.cert",
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	return exec.Command("sudo", "update-ca-certificates").Run()
}

func (l LinuxInstaller) Uninstall(certPath string) error {
	cmd := exec.Command(
		"sudo",
		"rm",
		"/usr/local/share/ca-certificates/myca.cert",
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
		"-addscore",
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
