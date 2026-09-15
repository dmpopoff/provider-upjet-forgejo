Patched svalabs/terraform-provider-forgejo v1.6.0: org/repo/team `id` is string
(FormatInt) so upjet IdentifierFromProvider works (empty id before Create);
Read/Delete 404 → gone/success for Crossplane finalizers.

Rebuild (linux/amd64):

```bash
cd providers/provider-forgejo
./hack/build-patched-tf-provider.sh
```

Full pipeline: docs/dev/upjet-provider-forgejo-build.md
