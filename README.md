# inscription-mint-bot
铭文批量操作工具
# 相关操作
## 执行命令
下载golang 安装包,本项目基于golang 1.20版本开发，建议不低于该版本
下载地址：https://go.dev/ 根据系统下相应的版本

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

进入相关目录下操作,例如eth系列的链进去evm目录执行下面命令
go mod tidy
go build -o eth-inscription-mint-bot //建议输出格式 xx-inscription-mint-bot xx替换成请求的链 比如 eth链为 eth-inscription-mint-bot avax链为 avax-inscription-mint-bot

## golang相关
安装golang后,建议配置GOPROXY,默认地址可能拉取依赖很慢

go env -w GOPROXY=https://goproxy.io,direct
go env -w GO111MODULE=on 该选项新版本默认都是开启的,根据go env 查看是否为on,不会查看的话直接执行就好了不会有影响
## (目录)/etc/deploy.yaml 修改

1.修改参数可以查看配置文件注释修改
# 备注
## 项目开发进度
1. 目前只开发了evm系列
## rpc接口参考地址
https://chainlist.org/chain/XX xx换成对应链id

## 题外话
1. 大佬的话可以打赏支持下,不强求奥,以后也会是开源
   打赏地址 0x04001842338fe79743680d8F3749eA53d16a41D9

2. 另外在avax链mint了个avi的铭文,有兴趣的小伙伴可以mint下
   https://avascriptions.com/token/detail?tick=avi

3. twitter地址： https://twitter.com/guq43432217 微信铭文技术交流群：
![img.png](img.png)
4. 欢迎有问题的小伙伴github上提Issue 也可以twitter,微信聊

最后满意的小伙伴不介意可以右上角点个 star