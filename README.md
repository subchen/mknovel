[![Build Status](https://travis-ci.org/subchen/mknovel.svg?branch=master)](https://travis-ci.org/subchen/mknovel)
[![License](http://img.shields.io/badge/License-Apache_2-red.svg?style=flat)](http://www.apache.org/licenses/LICENSE-2.0)


# mknovel

A general tool for download novel from website.

```
NAME:
   mknovel - Download a novel from URL, transform HTML to TEXT, pack it

USAGE:
   mknovel [options] URL

VERSION:
   1.2.4-47

AUTHORS:
   Guoqiang Chen <subchen@gmail.com>

OPTIONS:
   --threads num             parallel threads (default: 100)
   --short-chapter size      ignore chapter if size is short (default: 3000)
   -d dir, --directory dir   output directory (default: .)
   --help                    print this usage
   --version                 print version information
```

## Downloads

v1.2.4 Release: https://github.com/subchen/mknovel/releases/tag/v1.2.4

- Linux

    ```
    curl -fSL https://github.com/subchen/mknovel/releases/download/v1.2.4/mknovel-1.2.4-linux-amd64 -o /usr/local/bin/mknovel
    chmod +x /usr/local/bin/mknovel
    ```

- macOS

    ```
    curl -fSL https://github.com/subchen/mknovel/releases/download/v1.2.4/mknovel-1.2.4-darwin-amd64 -o /usr/local/bin/mknovel
    chmod +x /usr/local/bin/mknovel
    ```

- Windows

    ```
    wget https://github.com/subchen/mknovel/releases/download/v1.2.4/mknovel-1.2.4-windows-amd64.exe
    ```

## Usage

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
mknovel http://www.86696.cc/html/0/846/index.html
```
