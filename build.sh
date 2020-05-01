############################################
#当前版本号,每次更新服务时都必须更新版本号
############################################
CurrentVersion=0.1.1

############################################
#输出的二进制文件的 名称
############################################
OutFileName=PSSocketServer

############################################
#项目名称（包名称）看自己的go.mod的名字
############################################
Project=go_socketserver

############################################
#编译目标平台 (windows,linux...)
############################################
targetGOOS=linux
targetGOARCH=amd64


#############################
#####二进制生成文件存放路径#####
#############################
binPath=/Users/zhangfeng/Documents/GoProject/Bin/



















#########################################################################
####################### 以下如果没有需要无需改动  ##########################
#########################################################################
GitCommit=$(git rev-parse --short HEAD || echo unsupported)

fileName=""

if [ $targetGOOS = windows ]; then
  fileName=$OutFileName".exe"
else
   fileName=$OutFileName
fi
echo $fileName
GOOS=$targetGOOS GOARCH=$targetGOARCH go build -o $binPath'/'$OutFileName'/'$fileName  -ldflags "-X $Project/units.Version=$CurrentVersion -X '$Project/units.BuildTime=`date "+%Y-%m-%d %H:%M:%S"`' -X '$Project/units.GoVersion=`go version`' -X $Project/units.GitCommit=$GitCommit"

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

