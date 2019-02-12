# orator

## the simple static site generator

orator is a static site generator written in golang. it is very fast, easy to use
and flexible. orator takes a directory with content and renders it to html with
go templates called "layouts".

layouts and content files can include a yaml front matter that can be used in
the content. additionally, orator reads site-specific variables from a
`config.yaml` file. variables defined there can be accessed in all layouts and
content files.

orator runs on any platform where the go tool chain can run like plan 9, linux, windows,
mac os x and {dragonfly,free,open,net}bsd.

### installation

first, make sure that you have [go](https://golang.org) installed.

install orator with the following command:

```bash
go get github.com/tudurom/orator
```

### using orator

the first thing you need to do is setting up the directory structure.
an orator website needs a `config.yaml` file that stores site-wide configuration
and two directories: `layouts` for templates and `content` for your content.

#### creating a layout

**read more about go templates [here](https://golang.org/pkg/text/template/)**.

a simple layout may look like this:

`default.html`

```html
{{ define "head" }}
<head>
	<meta charset="utf-8">
	{{ $title := index .page.frontmatter "title" }}
	<title>{{ .siteconfig.title }}{{ if $title }} - {{ $title }}{{end}}</title>
</head>
{{ end }}

{{ define "header" }}
<header>
	<h1><a href="/">{{ .siteconfig.title }}</a></h1>
</header>
{{ end }}

{{ define "default" }}
<!doctype html>
{{ template "head" . }}
<body>
	{{ template "header" . }}
	{{ .page.content }}
</body>
{{ end }}
```

as you can see each file can contain multiple templates. they are all loaded
anyway.

#### assets

site's assets are stores in the `static` directory. they are automatically
copied in the root of the generated site's folder.

#### creating content

next up, we write some content for our site in the `content` directory. the
directory layout here is preserved in the generated site.

content can be in any format. if the file's name ends in `.md`, it will be
rendered as markdown to html.

content files can have a yaml front matter:

```markdown

---
layout: default
special_thing: false
---

hey this a site!

{{ if index .page.frontmatter "special_thing" }}
	<h1>hidden header</h1>
{{ end }}

* [users](/users)

```

in the example above, the header will not be shown because the
`special_thing` variable is set to false.

`layout` is a special variable that tells orator what layout should this page
use.

#### generating the site

`cd` into the site folder and run `./orator`. the final site will be generated
in the `gen` directory.
