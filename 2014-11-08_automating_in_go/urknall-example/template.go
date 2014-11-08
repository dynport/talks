package main

import (
	"fmt"
	"time"

	"github.com/dynport/urknall"
)

type Template struct {
}

func (tpl *Template) Render(p urknall.Package) {
	p.AddCommands("update",
		Shell(cacheBreaker()),
		UpdatePackages(),
	)

	// docker
	d := &Docker{Version: "1.3.1", Autostart: true}
	p.AddTemplate("docker-"+d.Version, d)

	// ruby
	rb := &Ruby{Version: "2.1.4"}
	p.AddTemplate("ruby-"+rb.Version, rb)

	// redis
	rd := &Redis{Version: "2.8.17"}
	p.AddTemplate("redis-"+rd.Version, rd)

	// nginx
	nx := &Nginx{Version: "1.4.7"}
	p.AddTemplate("nginx-"+nx.Version, nx)

	// postgres
	pg := &Postgres{Version: "9.3.4"}
	p.AddTemplate("postgres-"+pg.Version, pg)

	// redis
	//rd := &Redis{Version: "2.8.17"}
	//p.AddTemplate("redis-"+rd.Version, rd)

	// profile
	p.AddTemplate("user", &user{RubyPath: rb.InstallDir(), DockerPath: d.InstallDir()})
}

type user struct {
	RubyPath   string `urknall:"required=true"`
	DockerPath string `urknall:"required=true"`
}

func (tpl *user) Render(p urknall.Package) {
	p.AddCommands("profile", WriteFile("/root/.profile", profile, "root", 0644))
}

func cacheBreaker() string {
	return fmt.Sprintf("# cache breaker " + time.Now().Truncate(24*time.Hour).Format("2006-01-02"))
}

const profile = `#!/bin/bash
export PATH={{ .DockerPath }}/bin:{{ .RubyPath }}/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games
`
