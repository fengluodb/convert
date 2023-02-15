# Convert--一个格式转换工具

目前已实现的转换功能

- [x] txt格式向Epub格式转换。
- [x] png, jpg, gif格式相互转换

## 下载与安装
### 从源码编译
```
git clone https://github.com/fengluodb/convert.git
sh build.sh
```

## 使用方法

### 文本类转换
```shell
text subcommand converts text file into another format

Usage:
  convert text [flags]

Flags:
  -f, --format string        the target format
  -h, --help                 help for text
  -o, --output string        the output dir
  -s, --source stringArray   source files
```

示例
```shell
./convert text -s assets/斗破苍穹.txt -f epub -o test
```

### 图片类转换
```text
image subcommand converts image file into another image format

Usage:
  convert image [flags]

Flags:
  -f, --format string      the format of output files (default ".")
  -h, --help               help for image
  -o, --output string      the output dir (default ".")
  -s, --source strings     the list of source files
      --srcdir string      the dir of output files, it doesn't work if source paramter has values
      --srcformat string   the format of src files, it works with srcdir paramter
```

**示例**
```shell
# 提取出gif中的所有图片
./convert image -s assets/gopher.gif -f png -o test

# 将一个文件夹中的图片合成为gif图
./convert image --srcdir assets --srcformat png -f gif -o test
或者
./convert image -s assets/gopher-0.png -s assets/gopher-1.png -s assets/gopher-2.png -s assets/gopher-3.png -f gif -o test

# 将png转换为jpg
./convert image -s assets/gopher-0.png -f jpg -o test
也可以
./convert image --srcdir assets --srcformat png -f jpg -o test
```

## TODO
- [ ] 优化文本类相互转换的代码，增加对mobi，awz3格式的支持
- [ ] 增加resize图片、压缩图片的功能
- [ ] 增加与pdf相关的功能
- [ ] 开发convert在线版