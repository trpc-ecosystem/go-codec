English | [中文](README.zh_CN.md)

# rawstring

The rawstring protocol is a simple tcp-based protocol.

Its data format is: string `"(a=b&c=d)\n"`.

The default ending is `"\n"`, and the response is `"result=0&...\n"`.
