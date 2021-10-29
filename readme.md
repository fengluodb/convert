# Convert--一个文本格式转换工具

目前已实现的转换功能

- txt格式的网络小说向Epub格式转换。

## 使用方法

### txt向epub转换

在Releases下载convert后，执行以下两个命令之一

1. ./convert a.txt into epub
   
     使用时将a.txt替换为自己要转换的文件
2. ./convert a.txt b.jpg into epub

     使用时将a.txt， b.jpg转换为自己要转换的文件。相比于第一个命令，该命令可以将b.jpg作为epub的封面。

> 注意：目前只支持utf-8编码的txt文本，图片值支持jpg格式的图片