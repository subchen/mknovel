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
   2.0.0-64

AUTHORS:
   Guoqiang Chen <subchen@gmail.com>

OPTIONS:
       --novel-name value              name of novel
       --novel-author value            author of novel
       --novel-cover-image value       cover image file or url
       --input-encoding value          encoding for input txt file (default: GBK)
       --threads num                   parallel threads (default: 100)
       --short-chapter-size size       skip chapter if size is short (default: 3000)
       --format value                  output file format (txt, zip, epub) (default: epub)
   -d, --directory dir                 output directory (default: .)
       --output-encoding value         encoding for output txt file (default: GBK)
       --zip-filename-encoding value   encoding for output file name in zip (default: GBK)
       --help                          print this usage
       --version                       print version information

```

## Downloads

v2.0.0 Release: https://github.com/subchen/mknovel/releases/tag/v1.2.4

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

### Download from website

1. Create a config for target website

filename: `www.86696.cc.yaml`

```yaml
website-charset: GBK

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

cover-image:
  begin: <div id="fmimg">
  end: </div>
  regexp: src="(\S+)"
  name-index: 1

chapter-index:
  begin: <!--列表内容开始-->
  end: <!--列表内容结束-->
  regexp: <a href="([^>]+)">([^<]+)</a>
  link-index: 1
  name-index: 2

chapter-content:
  begin: <div id="BookText">
  end: </div>
```

2. Generate an epub novel

```bash
mknovel http://www.86696.cc/html/0/846/index.html
```

### Import a local txt file

1. Prepare a txt file, split chapters with double blank lines. 

```
第一章 北灵院
　　烈日如炎，灼热的阳光从天空...
　　在那一片投射着被柳树枝叶切...


第二章 被踢出灵路的少年
　　苏凌他们望着高台上的那些西...
    “喂，牧哥，那是西院的红...


第三章 牧域
    ...
```

2. 

```bash
mknovel 大主宰.txt --novel-name 大主宰 --novel-author 天蚕土豆 --novel-cover-image http://tu.zxcs8.com/content/uploadfile/201707/f3cc1499602096.jpg
```