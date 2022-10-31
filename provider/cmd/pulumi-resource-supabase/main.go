//go:generate go run ./generate.go

package main

import (
	"github.com/LuxChanLu/pulumi-supabase/pkg/provider"
	"github.com/LuxChanLu/pulumi-supabase/pkg/version"
)

var providerName = "supabase"

func main() {
	provider.Serve(providerName, version.Version, pulumiSchema)
}
