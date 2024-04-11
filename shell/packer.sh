unity_path="/e/Program Files/Unity 2021.3.30f1/Editor/Unity.exe" ##Unity路径
client_chuBao="/d/Yzr3Client/Client" #对应工程路径
echo -e "工程路径：$client_chuBao"

log_path="./log/"

today=$(date +%Y%m%d)
hour=$(date +%H)
minute=$(date +%M)
time_data="$today-$hour-$minute"

bundle_res_path="/e/ResourceBundle/YZR3V2HD"

#上传
obs_path="/d/jenkins/obsutil/aws"

#discord

apk_path="./build/"
apk_name="yzr3hd"

# 1. 打包产品
function build_product()
{
	local shellName="${client_chuBao}/Build/shell/build/build_android_apk_new.sh"
	"$shellName" -netWork ""\
	                  -clientChubao "$client_chuBao" \
	                  -clientChubaoCsharp "$client_chuBao" \
    								-bundleResPath "$bundle_res_path" \
    								-logPath "$log_path" \
                    -timeData "$time_data" \
                    -unityPath "$unity_path" \
                    -apkPath "$apk_path" \
                    -apkName "$apk_name" \
                    -buildOptions "BuildOptions.None" \
                    -bundleCode "PlayerSetting->Bundle Version Code" \
                    -buildVersion "1.2004.59" \
                    -serverTag "cn-release315" \
                    -identifier "" \
                    -svnup false \
                    -buildPkg true \
                    -mtpPkg false \
                    -obsPath "$obs_path" \
                    -uploadRes false \
                    -isXianFeng false

  # shellcheck disable=SC2181
  if [ "$?" -eq 0 ]; then ## $? 是shell 上一条命令的返回值，如果执行成功，退出码是 0，如果失败，退出码是 非0
    	echo "Execute Success........"
	else
    	echo "Execute Failed........."
     	exit 1
	fi
}

function buildAPK()
{
	echo -e "\033[41;33m 开始打包 android \033[0m"

	echo -e "\033[32m 步骤1：打包产品 \033[0m"
	build_product

	echo -e "\033[41;33m 打包结束 android \033[0m"
}

buildAPK