# svalabs id-string overlays

Incomplete copies of `svalabs/terraform-provider-forgejo` `internal/provider`
sources. Applied by [`../build-patched-tf-provider.sh`](../build-patched-tf-provider.sh)
onto a full upstream checkout.

Every `*.go` file starts with `//go:build ignore` so this Crossplane module’s
`go test` / golangci-lint **typecheck** does not treat them as part of
`github.com/crossplane-contrib/provider-upjet-forgejo` (they reference helpers
and packages that exist only in the svalabs tree). The build script strips that
constraint when copying.
