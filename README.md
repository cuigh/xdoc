# xdoc

A document site based Markdown.

![Snapshot](xdoc.png)

## Usage

### Standalone

```bash
xdoc -d /docs
```

### Docker

First you need create a new docker image base on `cuigh/xdoc`, copy all your documents to `/docs` directory in image. Here is the sample **Dockerfile** file.

```docker
FROM cuigh/xdoc
COPY . /docs/
```

Build the image

```bash
docker build -t docs .
```

Start the container

```bash
docker run -it -p 8000:8000 docs
```

## Customize menus

**xdoc** build document menus according to filename by default, you can customize it by adding a `menu.xml` on document root directory.

```xml
<menu>
    <item name="Java" url="/java/index.md">
        <item name="Library" url="##">
            <item name="Front" url="/java/lib/front.md"/>
            <item name="Utils" url="/java/lib/utils.md"/>
        </item>
        <item name="Tool" url="##">
            <item name="lark-cli" url="/java/tool/lark-cli.md"/>
            <item name="jctl" url="/java/tool/jctl.md"/>
        </item>
    </item>
</menu>
```

> NOTE: **xdoc** will look for `README.md` or `index.md` as index page automatically if exists.