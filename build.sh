#当前版本号,每次更新服务时都必须更新版本号
CurrentVersion=0.1.1
#项目名
Project=PSSocketServer

#编译目标平台 (windows,linux...)
targetGOOS=windows
targetGOARCH=amd64

#二进制生成文件存放路径
binPath=/Users/zhangfeng/Documents/GoProject/Bin/



#######################以上为配置项#######################




GitCommit=$(git rev-parse --short HEAD || echo unsupported)

fileName=""

if [ $targetGOOS = windows ]; then
  fileName=$Project".exe"
else
   fileName=$Project
fi
echo $fileName
GOOS=$targetGOOS GOARCH=$targetGOARCH go build -o $binPath'/'$Project'/'$fileName  -ldflags "-X main.Version=$CurrentVersion -X 'main.BuildTime=`date "+%Y-%m-%d %H:%M:%S"`' -X 'main.GoVersion=`go version`' -X main.GitCommit=$GitCommit"

if [ $? -eq 0 ]; then
  echo "##############################"
  echo "build succeed !!"
  echo "Version:" $CurrentVersion
  echo "Git commit:" $GitCommit
  echo "Go version:" `go version`
  echo "##############################"
else
  echo "##############################"
  echo "ERROR!"
  echo "##############################"
fi

