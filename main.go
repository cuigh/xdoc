package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/ext/files"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/net/web/renderer/std"
	"github.com/cuigh/auxo/util/cast"
	"github.com/cuigh/xdoc/menu"
	"github.com/russross/blackfriday"
)

var docsDir string

func main() {
	config.BindFlags(app.Flags.Inner())
	config.SetEnvPrefix("")

	app.Name = "xdoc"
	app.Version = "0.1"
	app.Desc = "A document site based Markdown"
	app.Action = func(ctx *app.Context) {
		docsDir = config.GetString("xdoc.dir")
		if docsDir == "" {
			fmt.Println("xdoc: document folder must be specified, try 'xdoc -h' for more information")
			os.Exit(1)
		}
		menu.Init(docsDir)
		app.Run(server())
	}
	app.Flags.Register(flag.Help | flag.Version | flag.Config)
	app.Flags.String("xdoc.dir", "d", "", "the directory of documents")
	app.Start()
}

func server() *web.Server {
	dir := filepath.Dir(app.Path())

	s := web.Auto()
	s.Renderer = std.Must()

	// register static handlers
	s.File("/favicon.ico", filepath.Join(dir, "assets/favicon.ico"))
	s.Static("/css", filepath.Join(dir, "assets/css"))
	s.Static("/js", filepath.Join(dir, "assets/js"))
	s.Static("/fonts", filepath.Join(dir, "assets/fonts"))

	s.Get("/", mark)
	s.Get("/*doc", mark)
	return s
}

func mark(c web.Context) error {
	u, p := findFile(c.Request().RequestURI, "index.md", "README.md", "index.html")
	if files.NotExist(p) {
		return web.ErrNotFound
	}

	if !strings.HasSuffix(p, ".md") {
		return c.Content(p)
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	content := cast.BytesToString(blackfriday.Run(b))
	d := data.Map{
		"Version": app.Version,
		"Content": template.HTML(content),
		"Menu":    menu.Get().GetInfo(u),
	}
	return c.Render("layout", d)
}

func findFile(uri string, names ...string) (u string, p string) {
	p = filepath.Join(docsDir, uri)
	if files.IsDir(p) {
		for _, name := range names {
			path := filepath.Join(p, name)
			if files.Exist(path) {
				return filepath.Join(uri, name), path
			}
		}
	}
	return uri, p
}
