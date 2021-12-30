# 快速扫码服务端
----
用于**B/S架构**快速部署手机扫码功能
<p align="center">
  <a href="https://github.com/ynqjwsm/fast-qrcode/actions">
    <img src="https://img.shields.io/github/workflow/status/ynqjwsm/fast-qrcode/Go?style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/ynqjwsm/fast-qrcode">
    <img src="https://goreportcard.com/badge/github.com/ynqjwsm/fast-qrcode?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/ynqjwsm/fast-qrcode?style=flat-square">
  <a href="https://github.com/ynqjwsm/fast-qrcode/releases">
    <img src="https://img.shields.io/github/release/ynqjwsm/fast-qrcode/all.svg?style=flat-square">
  </a>
</p>

### 特点
- 二维码数据预生成，提升请求响应时延
- 内置缓存，极简部署

###安装
```bash
# 将文件上传至服务器
# 解压文件
tar -zxvf fast-qrcode_v1.0.0_linux_amd64.tar.gz -C /opt

# 创建配置文件夹
mkdir -p /etc/fast-qrcode/

# 移动配置文件
mv /opt/default.conf /etc/fast-qrcode/default.conf

# 配置文件可执行
chmod +x /opt/fast-qrcode
```

###创建Service文件
```bash
vim /etc/systemd/system/fast-qrcode.service

#文件内容
[Unit]
Description=Fast QRScan Service
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
ExecStart=/opt/fast-qrcode -c /etc/fast-qrcode/default.conf

[Install]
WantedBy=multi-user.target
```

> 启动服务：systemctl start fast-qrcode

> 停止服务：systemctl stop fast-qrcode

> 设置自启：systemctl enable fast-qrcode

###接口说明

#### 1. 测量接口
```/actuator/:metric```

metric取值
- cache - 监控内部缓存运行情况

```EntryCount``` the number of items currently in the cache.

```EvacuateCount``` is a metric indicating the number of times an eviction occurred.

```ExpiredCount``` is a metric indicating the number of times an expire occurred.

```HitCount``` is a metric that returns number of times a key was found in the cache.

```HitRate``` is the ratio of hits over lookups.

```LookupCount``` is a metric that returns the number of times a lookup for a given key occurred.

```MissCount``` is a metric that returns the number of times a miss occurred in the cache.

```OverwriteCount``` indicates the number of times entries have been overriden.

```r``` 返回值: ```0``` 正常，```-1``` 异常

```json
{
  "EntryCount": 272,                          
  "EvacuateCount": 0,
  "ExpiredCount": 10,
  "HitCount": 7156,
  "HitRate": 0.8909362549800797,
  "LookupCount": 8032,
  "MissCount": 876,
  "OverwriteCount": 819,
  "r": 0
}
```
- create - 创建扫码请求计数器
- lookup - 轮询查找扫码结果计数器
- submit - 提交扫码结果计数器
- notify - 一次扫码通知计数器

#### 2. 服务存活测试
```/ping```

用于测量服务存活状态，正常返回值：
```json
{
  "message": "pong"
}
```