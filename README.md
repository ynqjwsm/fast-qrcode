# 快速扫码服务端
----
用于B/S架构快速部署手机扫码功能

#### 安装
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

#### 创建Service文件
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