[English](README.md) | 中文

# rawstring

rawstring 协议是一种简单的基于 tcp 的调用协议

其数据格式为：字符串 `"(a=b&c=d)\n"`

结尾默认是 `"\n"`，响应回复 `"result=0&...\n"`
