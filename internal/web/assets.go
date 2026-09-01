package web

import "embed"

// content contains the complete browser-facing UI so Shorty can be deployed as
// a single binary without a separate templates or assets directory.
//
//go:embed components/*.html layouts/*.html pages/*.html static/css/*.css static/js/*.js
var content embed.FS
