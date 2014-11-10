# bitbucket.org/ulfurinn/cli

This package is a fork/rewrite of `github.com/codegangsta/cli` with extended shell completion functionality. It uses `bitbucket.org/ulfurinn/options` instead of the standard `flag` package, and mixing both in the same application is not recommended.

If you want to migrate from `github.com/codegangsta/cli`, you'll need to make minor changes in your command setup as there are differences in exposed types.
