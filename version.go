package main

// Version is set at build time via -ldflags "-X main.Version=x.y.z".
// Falls back to "dev" for local builds.
var Version = "dev"
