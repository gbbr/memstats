package web

import (
	"html/template"
	"strings"
)

func Template() (*template.Template, error) {
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
		<title>MemViewer</title>
		<style>
			{{template "stylesheet"}}
		</style>
	</head>
	<body>
		<script id="ms-viewer-template" type="template/text">
		<div class="group">
			<h3>General</h3>
			<div class="cell">
				Allocated and using: <%= Alloc %>
			</div>
			<div class="cell">
				Total + Freed: <%= TotalAlloc %>
			</div>
			<div class="cell">
				System: <%= Sys %>
			</div>
			<br />
			<div class="cell">
				Lookups: <%= Lookups %>
			</div>
			<div class="cell">
				Frees: <%= Frees %>
			</div>
			<div class="cell">
				mallocs: <%= Mallocs %>
			</div>
		</div>

		<div class="group">
			<h3>Heap</h3>
			<div class="cell">
				Allocated and using: <%= HeapAlloc %>
			</div>
			<div class="cell">
				System: <%= HeapSys %>
			</div>
			<div class="cell">
				Idle: <%= HeapIdle %>
			</div>
			<div class="cell">
				In use: <%= HeapInuse %>
			</div>
			<div class="cell">
				Released: <%= HeapReleased %>
			</div>
			<br />
			<div class="cell">
				Objects: <%= HeapObjects %>
			</div>
		</div>

		<div class="group">
			<h3>Low-level allocator statistics</h3>
			<!-- Low-level fixed-size structure allocator statistics.-->
			<!--	Inuse is bytes used now.-->
			<!--	Sys is bytes obtained from system.-->
			<div class="cell">
				Stack:
				<%= StackInuse %> of <%= StackSys %>
			</div>
			<div class="cell">
				MSpan:
				<%= MSpanInuse %> of <%= MSpanSys %>
			</div>
			<div class="cell">
				MCache: <%= MCacheInuse %> of <%= MCacheSys %>
			</div>
			<br />
			<div class="cell">
				BuckHashSys: <%= BuckHashSys %>
			</div>
			<div class="cell">
				GCSys: <%= GCSys %>
			</div>
			<div class="cell">
				Other: <%= OtherSys %>
			</div>
		</div>

		<div class="group">
			<h3>Garbage collector</h2>
			<div class="cell">
				Next run: <%= NextGC %>
			</div>
			<div class="cell">
				Last run: <%= LastGC %>
			</div>
			<div class="cell">
				Pause: <%= PauseTotalNs %>
			</div>
			<div class="cell">
				Runs: <%= NumGC %>
			</div>
			<br />
			<div class="cell">
				Enabled: <%= EnableGC %>
			</div>
			<div class="cell">
				Debug: <%= DebugGC %>
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
