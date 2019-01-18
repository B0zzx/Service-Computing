# CLI-selpg

## 项目描述

+ 使用GO语言开发Linux命令行实用程序。
+ 该程序允许用户指定来自从输入文本抽取的页的范围，文本可以来自文件或另一进程。

+ 程序开发参考于[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) 。

## 数据结构：

```go
type selpgArgs struct{
    startPage int
    endPage int
    inFilename string
    pageLen int
    pageType int
    printDest string
}
```

+ pageLen 表示每页的行数，可以被 “-l” 指令修改；
+ pageType 表示页的类型，l 为确定页行数的文本，f 为页数由ASCII码确定的换页字符定界的文本；

## 测试结果：

1. `$ selpg`：

   ```shell
   selpg: not enough arguments
   
   USAGE: selpg -sstart_page -eend_page [ -f | -llines_per_page ][ -ddest ] [ in_filename ]
   ```

2. `$ selpg -s1 -e2 -l5 input.txt`，将输入文本的第一页和第二页输出，每一页有5行：

   ```
   input file line 1
   input file line 2
   input file line 3
   input file line 4
   input file line 5
   input file line 6
   input file line 7
   input file line 8
   input file line 9
   input file line 10
   ```

3. `$ selpg -s1 -e2 < input.txt`，读取标准输入（标准输入被 shell / 内核重定向为来自 “input_file” 而不是显式命名的文件名参数）。

4. `selpg -s2 -e4 -l10 input.txt >output.txt`，将第 2 页到第 4 页写至标准输出（标准输出被 shell／内核重定向至 “output.txt”），屏幕不显示时使用`cat output.txt`命令在终端显示 output.txt 内容：

   ```
   input file line 11
   input file line 12
   input file line 13
   input file line 14
   input file line 15
   input file line 16
   input file line 17
   input file line 18
   input file line 19
   input file line 20
   input file line 21
   input file line 22
   input file line 23
   input file line 24
   input file line 25
   input file line 26
   input file line 27
   input file line 28
   input file line 29
   input file line 30
   input file line 31
   input file line 32
   input file line 33
   input file line 34
   input file line 35
   input file line 36
   input file line 37
   input file line 38
   input file line 39
   input file line 40
   ```

5. `selpg -s10 -e20 input.txt 2>error.txt`，不符合标准的信息将被输出至错误信息文件error.txt，此时，屏幕不显示；我们使用`cat error.txt`命令在终端显示内容：

   ```
   selpg: startPage (10) greater than total pages (1), no output written
   ```

6. `selpg -s1 -e2 -f input.txt`，假定页由换页符定界，第 1 页到第 2 页被写至 selpg 的标准输出；

7. `selpg -s1 -e2 -dlp1 input_file`，输出错误信息：

   ```
   selpg: could not open pipe to ""
   ```
