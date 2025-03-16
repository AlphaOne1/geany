<p align="center">
    <!--<img src="geany_logo.svg" width="25%" alt="Logo"><br>-->
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/test.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/test.yml/badge.svg"
             alt="Test Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/codeql.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/codeql.yml/badge.svg"
             alt="CodeQL Pipeline Result">
    </a>
    <a href="https://github.com/AlphaOne1/geany/actions/workflows/security.yml"
       rel="external"
       target="_blank">
        <img src="https://github.com/AlphaOne1/geany/actions/workflows/security.yml/badge.svg"
             alt="Security Pipeline Result">
    </a>
    <a href="https://goreportcard.com/report/github.com/AlphaOne1/geany"
       rel="external"
       target="_blank">
        <img src="https://goreportcard.com/badge/github.com/AlphaOne1/geany"
             alt="Go Report Card">
    </a>
    <a href="https://codecov.io/github/AlphaOne1/geany"
       rel="external"
       target="_blank">
        <img src="https://codecov.io/github/AlphaOne1/geany/graph/badge.svg?token=SIQ3UG8OJI"
             alt="Code Coverage">
    </a>
    <!--<a href="https://www.bestpractices.dev/projects/0000"
       rel="external"
       target="_blank">
        <img src="https://www.bestpractices.dev/projects/0000/badge"
             alt="OpenSSF Best Practises">
    </a>-->
    <a href="https://scorecard.dev/viewer/?uri=github.com/AlphaOne1/geany"
       rel="external"
       target="_blank">
        <img src="https://api.scorecard.dev/projects/github.com/AlphaOne1/geany/badge"
             alt="OpenSSF Scorecard">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany?ref=badge_shield&issueType=license"
       rel="external"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany.svg?type=shield&issueType=license"
            alt="FOSSA Status">
    </a>
    <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany?ref=badge_shield&issueType=security" 
       rel="external"
       target="_blank">
        <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FAlphaOne1%2Fgeany.svg?type=shield&issueType=security"
             alt="FOSSA Status">
    </a>
    <!--<a href="https://godoc.org/github.com/AlphaOne1/geany"
       rel="external"
       target="_blank">
        <img src="https://godoc.org/github.com/AlphaOne1/geany?status.svg"
             alt="GoDoc Reference">
    </a>-->
</p>

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
    \ \____.----(O)----..---.
     `-__  o \_______/o \ O /
         `\  O O O O O  /`-`
           `-=========-`
```

geany
=====

*geany* is a library to easily print logos enriched with custom information using the Go templating engine. 

Installation
------------

To install *geany*, you can use the following command:

```bash
$ go get github.com/AlphaOne1/geany
```

Getting Started
---------------

### Simple Case

Assuming a file `logo.tmpl` as this one:

```text
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
    \ \____.----(O)----..---.
     `-__  o \_______/o \ O /
         `\  O O O O O  /`-`
           `-=========-`
```

a simple program to produce the logo could look like this:

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
  \  /                on unknown
   \ \     from revision unknown 
   / /
   |/
   (o)
    \ \____.----(O)----..---.
     `-__  o \_______/o \ O /
         `\  O O O O O  /`-`
           `-=========-`
```

### User Provided Data

Assuming that the *geany* supplied information is not enough, a user can provide additional data. This data
is then accessible inside of the template using `.Value`. A modified example template could look like this:

```text
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
    \ \____.----(O)----..---.
     `-__  o \_______/o \ O /
         `\  O O O O O  /`-`
           `-=========-`
```

The Go program of the [simple case](#simple-case) has only to be slightly modifed, providing the new data.

```go
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

//go:embed logo.tmpl
var logo string

func main() {
	_ = geany.PrintLogo(logo, &struct {
		FeatureA bool
		FeatureB bool
        Greeting string
	}{FeatureA: true, FeatureB: false, Greeting: "Hi Geany!"})
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
  \  /                on unknown
   \ \     from revision unknown 
   / /
   |/   Special Features A: enabled
   (o)                   B: disabled
    \ \____.----(O)----..---.
     `-__  o \_______/o \ O /
         `\  O O O O O  /`-`
           `-=========-`
```

Other than using a `struct`, also a `map` could be used to introduce the data into the template. The keys of the map
then take the place of the structure members.

All examples can be found in [examples](examples) folder.