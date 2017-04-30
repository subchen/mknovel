[![License](http://img.shields.io/badge/License-Apache_2-red.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)


# mknovel

A general tool for download novel from website.

```
Usage: mknovel [--threads=100] [--short-chapter=3000] [-d dir] URL
   or: mknovel [ --version | --help ]

Download a novel from URL, transform HTML to TEXT, zipped it.

Options:
  --threads=100          parallel threads
  --short-chapter=3000   ignore short chapter
  -d, --directory=.      output directory
  --version              show version information
  --help                 show this help
```

# Downloads

Release: `1.2.1`

* Linux-amd64: https://raw.githubusercontent.com/subchen/mknovel/master/bin/mknovel-linux-1.2.1
* MacOS-amd64: https://raw.githubusercontent.com/subchen/mknovel/master/bin/mknovel-darwin-1.2.1
* Window-amd64: https://raw.githubusercontent.com/subchen/mknovel/master/bin/mknovel-windows-1.2.1.exe

# Build from Source

```bash
# install glide
curl https://glide.sh/get | sh

glide install
make build
```

# Usage

1. Create a config for target website

filename: `www.86696.cc.yaml`

```yaml
website-charset: GBK
zipfilename-charset: GBK

title:
    begin: <title>
    end: </title>
    regexp: (.+)最新章节
    name-index: 1

author:
    begin: <div id="info">作者：<a
    end: </a>
    regexp: href=".+" target="_blank">(.+)
    name-index: 1

chapter:
    begin: <!--列表内容开始-->
    end: <!--列表内容结束-->
    regexp: <a href="([^>]+)">([^<]+)</a>
    link-index: 1
    name-index: 2

content:
    begin: <div id="BookText">
    end: </div>
```

2. Build novel zip

```bash
mknovel http://www.86696.cc/html/0/846/index.html -d "."
```
