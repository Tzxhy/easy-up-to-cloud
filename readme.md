# 构建

windows: CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o easy_up_cloud.exe
linux: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o easy_up_cloud_linux
mac: CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o easy_up_cloud_mac