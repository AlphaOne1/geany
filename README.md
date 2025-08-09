<!-- markdownlint-disable MD013 MD033 MD041 -->
<p align="center">
    <!--<img src="geany_logo.svg" width="25%" alt="Logo"><br>-->
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/test.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/test.yml/badge.svg"
             alt="Test Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/codeql.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/codeql.yml/badge.svg"
             alt="CodeQL Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/security.yml"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/security.yml/badge.svg"
             alt="Security Pipeline Result">
    </a>
    <a href="https://goreportcard.com/report/github.com/AlphaOne1/geany"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://goreportcard.com/badge/github.com/AlphaOne1/geany"
             alt="Go Report Card">
    </a>
    <a href="https://codecov.io/gh/AlphaOne1/geany"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://codecov.io/gh/AlphaOne1/geany/graph/badge.svg"
             alt="Code Coverage">
    </a>
    <a href="https://coderabbit.ai"
       rel="external noopener noreferrer"
       target="_blank">
       <img src="https://img.shields.io/coderabbit/prs/github/AlphaOne1/geany"
            alt="CodeRabbit Reviews">
    </a>
    <!--<a href="https://www.bestpractices.dev/projects/0000"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://www.bestpractices.dev/projects/0000/badge"
             alt="OpenSSF Best Practises">
    </a>-->
    <a href="https://scorecard.dev/viewer/?uri=github.com/AlphaOne1/geany"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://api.scorecard.dev/projects/github.com/AlphaOne1/geany/badge"
             alt="OpenSSF Scorecard">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany?ref=badge_shield&issueType=license"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany.svg?type=shield&issueType=license"
            alt="FOSSA Status">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany?ref=badge_shield&issueType=security"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany.svg?type=shield&issueType=security"
             alt="FOSSA Status">
    </a>
    <a href="https://pkg.go.dev/github.com/AlphaOne1/geany"
       rel="external noopener noreferrer"
       target="_blank">
        <img src="https://pkg.go.dev/badge/github.com/AlphaOne1/geany.svg"
             alt="GoDoc Reference">
    </a>
</p>
<!-- markdownlint-enable MD013 MD033 MD041 -->

```text
   ..
  =≙≙=
 _.OO._
/ \__/ \
 \_><_/
 |=%%=|
  \  /
   \ \
   / /
   |/
   (o)
    \ \____.----(O)----.,.---.
     `-__  o \_______/ o | O |
         `\  O O O O O  / `-`
           `-=========-`
```

geany
=====

*geany* is a library to easily print logos enriched with custom information using the Go templating engine.

Writing a program and just wanting a nice logo for the start should be without hassles. *geany* aims to
provide a smooth experience for developers that just want to show a nice logo and want to enrich it with
some more information.


Installation
------------

To install *geany*, you can use the following command:

```bash
$ go get github.com/AlphaOne1/geany
```


Getting Started
---------------

### Simple Case

Assuming a file `logo.tmpl` like this one:

```gotemplate
   ..
  =≙≙=
 _.OO._
/ \__/ \
 \_><_/
 |=%%=|      Build using {{ .Geany.GoVersion   }}
  \  /                on {{ .Geany.VcsTime     }}
   \ \     from revision {{ .Geany.VcsRevision }} {{ if eq .Geany.VcsModified "*" }} (modified) {{ end }}
   / /
   |/
   (o)
    \ \____.----(O)----.,.---.
     `-__  o \_______/ o | O |
         `\  O O O O O  / `-`
           `-=========-`
```

A simple program to produce the logo could look like this:

```go
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

//go:embed logo.tmpl
var logo string

func main() {
	_ = geany.PrintLogo(logo, nil)
}
```

This program produces the following output:

```text
   ..
  =≙≙=
 _.OO._
/ \__/ \
 \_><_/
 |=%%=|      Build using go1.24.1
  \  /                on 2025-03-16T03:16:43Z
   \ \     from revision ce78d813448120739f345efa679b2b244a42e679
   / /
   |/
   (o)
    \ \____.----(O)----.,.---.
     `-__  o \_______/ o | O |
         `\  O O O O O  / `-`
           `-=========-`
```

*geany* provides these attributes inside `.Geany`

| Name        | Description                                                                                               |
|-------------|-----------------------------------------------------------------------------------------------------------|
| GoVersion   | Go Version used to build the executable                                                                   |
| VcsModified | Contains a `*` if the content of the repository was modified before the build, if version control is used |
| VcsRevision | Revision of version control system, if used                                                               |
| VcsTime     | Time of the revision in the version control system, if used                                               |


### User Provided Data

Assuming that the *geany* supplied information is not enough, a user can provide additional data. This data
becomes accessible in the template using `.Value`. A modified example template could look like this:

```gotemplate
   ..      ________________________
  =≙≙=    / {{ .Values.Greeting }}
 _.OO._  /
/ \__/ \
 \_><_/
 |=%%=|      Build using {{ .Geany.GoVersion   }}
  \  /                on {{ .Geany.VcsTime     }}
   \ \     from revision {{ .Geany.VcsRevision }} {{ if eq .Geany.VcsModified "*" }} (modified) {{ end }}
   / /
   |/   Special Features A: {{ if .Values.FeatureA -}} enabled {{- else -}} disabled {{- end }}
   (o)                   B: {{ if .Values.FeatureB -}} enabled {{- else -}} disabled {{- end }}
    \ \____.----(O)----.,.---.
     `-__  o \_______/ o | O |
         `\  O O O O O  / `-`
           `-=========-`
```

The Go program of the [simple case](#simple-case) needs only slight modification to provide the new data:

```go
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

//go:embed logo.tmpl
var logo string

func main() {
	_ = geany.PrintLogo(logo,
		&struct {
			FeatureA bool
			FeatureB bool
			Greeting string
		}{
			FeatureA: true,
			FeatureB: false,
			Greeting: "Hi Geany!",
		})
}
```

This modified version would print now:

```text
   ..      ________________________
  =≙≙=    / Hi Geany!
 _.OO._  /
/ \__/ \
 \_><_/
 |=%%=|      Build using go1.24.1
  \  /                on 2025-03-16T03:16:43Z
   \ \     from revision ce78d813448120739f345efa679b2b244a42e679
   / /
   |/   Special Features A: enabled
   (o)                   B: disabled
    \ \____.----(O)----.,.---.
     `-__  o \_______/ o | O |
         `\  O O O O O  / `-`
           `-=========-`
```

Instead of using a `struct`, a `map` could be used to pass the data into the template. The keys of the map
then take the place of the structure members.


### Basic or Fallback Output

*geany* offers a fallback mechanism, should a problem occur in normal logo output. Should normal operation fail, the
`PrintSimple` function is employed to somehow give some information, without having the normal logo around.

Note that the logo template itself is checked always and the program will report a broken logo template visibly. So this
is __not__ intended to cover up poor logo template design.

The following program uses `PrintSimple` directly:

```go
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

func main() {
	_ = geany.PrintSimple(nil)
}
```

It produces the following output:

```text
./logo_simple
{
    "Geany": {
        "VcsRevision": "5f4b130546ef200186692d4f88c83386f0c4ae98",
        "VcsTime": "2025-03-16T03:21:23Z",
        "VcsModified": "*",
        "GoVersion": "go1.24.1"
    },
    "Values": null
}
```

All examples can be found in [examples](examples) folder.