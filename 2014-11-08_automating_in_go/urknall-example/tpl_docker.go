package main

import "github.com/dynport/urknall"

type Docker struct {
	Version          string `urknall:"required=true"` // e.g. 1.1.0
	CustomInstallDir string
	Public           bool
	Autostart        bool
}

func (docker *Docker) Render(pkg urknall.Package) {
	pkg.AddCommands("packages", InstallPackages("aufs-tools", "cgroup-lite", "xz-utils", "git"))
	pkg.AddCommands("install",
		Mkdir("{{ .InstallDir }}/bin", "root", 0755),
		Download("http://get.docker.io/builds/Linux/x86_64/docker-{{ .Version }}", "{{ .InstallDir }}/bin/docker", "root", 0755),
	)
	pkg.AddCommands("upstart", WriteFile("/etc/init/docker.conf", dockerUpstart, "root", 0644))
	if docker.Autostart {
		pkg.AddCommands("start", Shell("if status docker | grep running; then restart docker; else start docker; fi"))
	}
}

const dockerUpstart = `exec {{ .InstallDir }}/bin/docker -d -H tcp://{{ if .Public }}0.0.0.0{{ else }}127.0.0.1{{ end }}:4243 -H unix:///var/run/docker.sock 2>&1 | logger -i -t docker
`

func (docker *Docker) InstallDir() string {
	if docker.Version == "" {
		panic("Version must be set")
	}
	if docker.CustomInstallDir != "" {
		return docker.CustomInstallDir
	}
	return "/opt/docker-" + docker.Version
}
