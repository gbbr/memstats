package view

import (
	"html/template"
	"strings"
)

func Render() (*template.Template, error) {
	return template.New("name").Parse(strings.Join([]string{
		underscoreJS,
		mainJS,
		stylesheet,
		rootView,
	}, ""))
}

var rootView = `
{{define "main"}}
<!DOCTYPE html>
<html>
	<head>
		<title>MemStats Viewer</title>
		<style>
			{{template "stylesheet"}}
		</style>
	</head>
	<body>
		<script id="ms-viewer-template" type="template/text">
		<div class="group">
			<h3>General</h3>
			<div class="cell">
				Allocated and using: <%= Stats.Alloc %>
			</div>
			<div class="cell">
				Total + Freed: <%= Stats.TotalAlloc %> 
			</div>
			<div class="cell">
				System: <%= Stats.Sys %> 
			</div>
			<br />
			<div class="cell">
				Lookups: <%= Stats.Lookups %>
			</div>
			<div class="cell">
				Frees: <%= Stats.Frees %>
			</div>
			<div class="cell">
				mallocs: <%= Stats.Mallocs %>
			</div>
		</div>

		<div class="group">
			<h3>Heap</h3>
			<div class="cell">
				Allocated and using: <%= Stats.HeapAlloc %> 
			</div>
			<div class="cell">
				System: <%= Stats.HeapSys %> 
			</div>
			<div class="cell">
				Idle: <%= Stats.HeapIdle %> 
			</div>
			<div class="cell">
				In use: <%= Stats.HeapInuse %> 
			</div>
			<div class="cell">
				Released: <%= Stats.HeapReleased %> 
			</div>
			<br />
			<div class="cell">
				Objects: <%= Stats.HeapObjects %>
			</div>
		</div>

		<div class="group">
			<h3>Low-level allocator statistics</h3>
			<!-- Low-level fixed-size structure allocator statistics.-->
			<!--	Inuse is bytes used now.-->
			<!--	Sys is bytes obtained from system.-->
			<div class="cell">
				Stack:
				<%= Stats.StackInuse %> of <%= Stats.StackSys %> 
			</div>
			<div class="cell">
				MSpan:
				<%= Stats.MSpanInuse %> of <%= Stats.MSpanSys %> 
			</div>
			<div class="cell">
				MCache: <%= Stats.MCacheInuse %> of <%= Stats.MCacheSys %> 
			</div>
			<br />
			<div class="cell">
				BuckHashSys: <%= Stats.BuckHashSys %>
			</div>
			<div class="cell">
				GCSys: <%= Stats.GCSys %>
			</div>
			<div class="cell">
				Other: <%= Stats.OtherSys %>
			</div>
		</div>

		<div class="group">
			<h3>Garbage collector</h2>
			<div class="cell">
				Next run: <%= Stats.NextGC %>
			</div>
			<div class="cell">
				Last run: <%= Stats.LastGC %>
			</div>
			<div class="cell">
				Pause: <%= Stats.PauseTotalNs %>
			</div>
			<div class="cell">
				Runs: <%= Stats.NumGC %>
			</div>
			<br />
			<div class="cell">
				Enabled: <%= Stats.EnableGC %>
			</div>
			<div class="cell">
				Debug: <%= Stats.DebugGC %>
			</div>
		</div>

		</script>
		<div id="ms-viewer"></div>

		<script>{{template "underscoreJS"}}</script>
		<script>{{template "mainJS"}}</script>
	</body>
</html>
{{end}}
`
