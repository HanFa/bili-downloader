
# BDownloader v0.1

This downloader is based on the [bili-go](https://github.com/HanFa/bili-go), the collection of Bilibili APIs.

Please do not use this repo for any monetized purposes.

BDownloader v0.1 is a light-weight cross-platform downloader for bilibili.tv. The default video output path
is `$HOME/<bvid>.mp4`.

![bdownloader](./doc/bdownloader.png)

Prereqs:
* golang 1.16

Usage:
```shell
$ make clean && make && ./output/bdownloader
```
