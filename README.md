# git-kustomize-diff

[![Go Reference][godoc-image]][godoc-link]
[![Coverage Status][cov-image]][cov-link]

Check diff of your kustomize directory.

![](misc/example.png)

## Prerequisites

- [golangci-lint v1.42.1](https://github.com/golangci/golangci-lint)

## Installation

Install from the Go source code:

```bash
$ go install -u github.com/dtaniwaki/git-kustomize-diff
```

For MacOS, use Homebrew:

```bash
brew tap dtaniwaki/git-kustomize-diff
brew install git-kustomize-diff
```

Or, download the binary from [GitHub Releases](https://github.com/dtaniwaki/git-kustomize-diff/releases) and put it in your `$PATH`.

## Usage

```bash
$ git-kustomize-diff run
```

Flags:

```
Usage:
  git-kustomize-diff run target_dir [flags]

Flags:
      --allow-dirty                        allow dirty tree
      --base string                        base commitish (default to origin/main)
      --debug                              debug mode
      --exclude string                     exclude regexp (default to none)
      --git-path string                    path of a git binary (default to git)
  -h, --help                               help for run
      --include string                     include regexp (default to all)
      --kustomize-load-restrictor string   kustomize load restrictor type (default to kustomizaton provider defaults)
      --kustomize-path string              path of a kustomize binary (default to embedded)
      --target string                      target commitish (default to the current branch)
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new [Pull Request](../../pull/new/master)

## Copyright

Copyright (c) 2021 Daisuke Taniwaki. See [LICENSE](LICENSE) for details.


[godoc-image]: https://pkg.go.dev/badge/github.com/dtaniwaki/git-kustomize-diff.svg
[godoc-link]: https://pkg.go.dev/github.com/dtaniwaki/git-kustomize-diff
[cov-image]:   https://coveralls.io/repos/github/dtaniwaki/git-kustomize-diff/badge.svg?branch=main
[cov-link]:    https://coveralls.io/github/dtaniwaki/git-kustomize-diff?branch=main

