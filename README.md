[![License](http://img.shields.io/badge/License-Apache_2-red.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)


# mknovel

A general tool for download novel from website.

```
Usage: mknovel [-d dir] URL
   or: mknovel [ --version | --help ]

Download a novel from URL, transform HTML to TEXT, zipped it.

Options:
  -d, --directory=.   download novel into directory
  --version           show version information
  --help              show this help
```

# Downloads

Release: `1.0.0`

* Linux-amd64: https://github.com/subchen/mknovel/bin/mknovel-linux-1.0.0
* MacOS-amd64: https://github.com/subchen/mknovel/bin/mknovel-darwin-1.0.0
* Window-amd64: https://github.com/subchen/mknovel/bin/mknovel-1.0.0.exe

# Usage

## Create a config for target website

filename: `www.86696.cc.yaml`

```yaml
website-charset: GBK
zipfilename-charset: GBK

title:
    begin: <title>
    end: </title>
    regexp: (.+)最新章节
    name-index: 1

chapter:
    begin: <!--列表内容开始-->
    end: <!--列表内容结束-->
    regexp: <a href="(.+)">(.+)</a>
    link-index: 1
    name-index: 2

content:
    begin: <div id="BookText">
    end: </div>
```

## Build novel zip

```bash
mknovel http://www.86696.cc/html/0/846/index.html -d "."
```
