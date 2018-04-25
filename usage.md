# Usage

## 自动启动

支持在 `session` 登录时自动启动需要的命令，需要启动的命令要写在定义好的配置文件中，配置文件信息如下：

+ 路径

    `/usr/share/multi-display-session/autostart` 或 `$HOME/.config/multi-display-session/autostart`

+ 内容

    文件中一行代表一个命令，必须要为绝对路径，参数加在后面，与在终端执行时的一样。

    如启动 `/usr/local/bin/open_terminal.sh`, 就将 `/usr/local/bin/open_terminal.sh` 写入到文件。

    需要启动多个命令就写多行


## 显示设置

可以在文件 `/usr/share/multi-display-session/outputs.json` 或 `$HOME/.config/multi-display-session/outputs.json` 中配置多屏模式，文件格式如下：

``` json
[
    {
        "Name": "xxx",
        "X": 0,
        "Y": 0,
        "Width": 1366,
        "Height": 768
    },
    {
        "Name": "xxx1",
        "X": 1366,
        "Y": 0,
        "Width": 1366,
        "Height": 768
    }
]
```


## 快捷键

`session` 预定义了一些快捷键，如下：

快捷键 | 动作
 ----  | ----
 Win-1 | 打开终端
 Win-2 | 打开终端
 Win-3 | 打开终端
 Win-r | 启动调试 `session`
 Win-t | 打开终端
 Win-delete | 注销
