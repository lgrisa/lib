#外部传入参数
isSvnUp=false #是否svn up
ServerTag="cn-ios315"
BuildVersion="1.2004.57"
BuildValue="1"
Identifier="1"
BuildPkg=true #是否出包
UploadRes=false #是否上传资源
OutPutPath="/Users/Shared/buildXCode"
XcodeName="yzr3"
buildIpa=true

#进入到脚本路径启动

#系统配置参数
unity_path="/Applications/Unity/Unity.app/Contents/MacOS/Unity"
log_path="/Users/Shared/product/IOS/IOSLog/"
obs_path="/Users/admin/Tools/obsutil"
bundle_res_path="/Users/Shared/ResourceBundle"

clientPath=$(pwd)/../../../
echo -e "clientPath =========== $clientPath ==========="

#上传
upload_res="$UploadRes"
obs_path="/Users/admin/Tools/obsutil"

# 1. 打包产品
function build_product()
{
local shellName="${clientPath}/Build/shell/build/build_ipa_new.sh"
echo "shellName========== ""$shellName"
bash "$shellName" -netWork "" \
                  -clientChubao "$clientPath" \
                  -bundleResPath "$bundle_res_path" \
                  -logPath "$log_path" \
                  -timeData "$(date +%Y%m%d)-$(date +%H)-$(date +%M)" \
                  -unityPath "$unity_path" \
                  -outputPath "$OutPutPath" \
                  -xcodeName "$XcodeName" \
                  -buildOptions "BuildOptions.None" \
                  -buildValue "$BuildValue" \
                  -buildVersion "$BuildVersion" \
                  -serverTag "$ServerTag" \
                  -identifier "$Identifier" \
                  -svnup "$isSvnUp" \
                  -buildPkg "$BuildPkg" \
                  -obsPath "$obs_path" \
                  -uploadRes "$upload_res" \
                  -buildIpa "$buildIpa"
}


function buildIPA()
{
echo -e "\033[41;33m 开始打包\033[0m"

echo -e "\033[32m 步骤1：打包产品 \033[0m"
build_product

echo -e "\033[41;33m 打包结束 \033[0m"
}

buildIPA