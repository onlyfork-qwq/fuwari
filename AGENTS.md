# 提交规范

1. Git 提交信息（commit message）必须使用简体中文。
2. 推送到远端仓库前，必须启用系统代理（本机 `127.0.0.1:10808`）。
3. 发布新文章时，提交信息必须严格使用以下格式：
   `posts:发布新文章《XXX》。XXXXX。`

# 参考命令

```powershell
git commit -m "feat: 使用简体中文提交信息"
```

```powershell
$env:HTTP_PROXY="http://127.0.0.1:10808"
$env:HTTPS_PROXY="http://127.0.0.1:10808"
$env:ALL_PROXY="socks5://127.0.0.1:10808"
git push origin main
```
