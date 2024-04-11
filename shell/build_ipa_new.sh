#!/bin/sh
net_work=""				 #inner outer
client_chubao=""
bundle_res_path=""
log_path=""
unity_path=""
build_version=""
server_tag=""
identifier=""
B_build_bundle=""
B_build_pkg=""
B_onlybuild_bundle=""

platform="iOS"

output_path=""
xcode_name=""

#工程参数
build_options=""
build_value=""

#签名
keystore_name=""
keystore_pass=""
keyalias_name=""
keyalias_pass=""

#上传
B_upload_res=""
obs_path=""
cdn_path=""

svnpath="/usr/local/Cellar/subversion/1.14.2_6/bin/svn"

#输出结果
build_successful="Build Successful!"
build_failed="Build Failed!"
build_result=""
color_green=65280
color_red=10027008
color_result=""

B_Build_Ipa=false

#是否为一个工程
b_isoneproj=false

xcode_proj_path=""

while [[ -n "$1" ]]; do
	case "$1" in
		-netWork )
			net_work="$2"
			echo -e "\n===========net_work:$net_work"
			shift
			;;
		-clientChubao )
			client_chubao="$2"
			echo "===========client_chubao:$client_chubao"
			shift
			;;
		-bundleResPath )
			bundle_res_path="$2"
			echo "===========bundle_res_path:$bundle_res_path"
			shift
			;;
		-logPath )
			log_path="$2"
			echo "===========log_path:$log_path"
			shift
			;;
		-timeData )
			time_data="$2"
			echo "===========time_data:$time_data"
			shift
			;;
		-unityPath )
			unity_path="$2"
			echo "===========unity_path:$unity_path"
			shift
			;;
		-outputPath )
			output_path="$2"
			echo "===========output_path:$output_path"
			shift
			;;
		-xcodeName )
			xcode_name="$2"
			echo "===========xcode_name:$xcode_name"
			shift
			;;
		-buildOptions )
            build_options="$2"
			echo "===========build_options:$build_options"
			shift
			;;
		-buildValue )
            build_value="$2"
			echo "===========build_value:$build_value"
			shift
			;;
		-buildVersion )
            build_version="$2"
			echo "===========build_version:$build_version"
			shift
			;;
		-serverTag )
            server_tag="$2"
			echo "===========server_tag:$server_tag"
			shift
			;;
		-identifier )
            identifier="$2"
			echo "===========identifier:$identifier"
			shift
			;;
		-svnup )
            B_svnup="$2"
			echo "===========B_svnup:$B_svnup"
			shift
			;;
		-buildBundle )
            B_build_bundle="$2"
			echo "===========B_build_bundle:$B_build_bundle"
			shift
			;;
		-onlyBuildAssets )
            B_onlybuild_bundle="$2"
			echo "===========B_onlybuild_bundle:$B_onlybuild_bundle"
			shift
			;;
		-buildPkg )
            B_build_pkg="$2"
			echo "===========B_build_pkg:$B_build_pkg"
			shift
			;;
		-obsPath )
            obs_path="$2"
			echo "===========obs_path:$obs_path"
			shift
			;;
		-uploadRes )
            B_upload_res="$2"
			echo "===========B_upload_res:$B_upload_res"
			shift
			;;
	  -buildIpa )
	          B_Build_Ipa="$2"
      echo "===========BuildIpa:$B_Build_Ipa"
      shift
      ;;
	esac
	shift
done


#Log颜色处理-------------start
infolog()
{
    echo -e Log: "\033[32m$1\033[0m"
}

#输出报错log
errorlog()
{
    echo -e Error: "\033[01;31m$1\033[0m"
}

#输出报错log并停止运行
errorlogExit()
{
    #echo -e Error: "\033[01;31m$1\033[0m"
    #echo -e Error: "\033[01;31m打包失败\033[0m"
	
	build_result="Build Failed : $1"
	color_result=$color_red
	#send_msg_to_discord
	
	
    echo -e Error: "$1"
    echo -e Error: "build failed!"
	exit 1
}
#Log颜色处理-------------end



arrange_keystore()
{
	keystore_name="${client_chubao_csharp}/Build/com.sgame.ex.keystore"
	keystore_pass="123456"
	keyalias_name=`basename $keystore_name`
	keyalias_pass="123456"
}

revert_version()
{
	cd ${client_path}"/ProjectSettings"
	svn revert -R * 
	svn up --username="chenli" --password="chenli"
}

svn_revert_chubao()
{
	echo -e "\n---------------svn_revert_chubao---------------\n"
	log_file="$log_path""$time_data""_svnup_chubao.log"
	if [[ ! -d "$log_path" ]]; then
		mkdir "$log_path"
	fi
	if [[ -f "$log_file" ]]; then
		rm -rf "$log_file"
	fi

	client_path=$1
	cd $client_path
	echo "cd "`pwd`
	$svnpath status
	$svnpath cleanup
	{
		# find . -not \( -path ./Assets/Resources -prune \) -not \( -path ./Library -prune \) -type f | xargs $svnpath revert
		# find . -not \( -path ./Assets/Resources -prune \) -not \( -path ./Library -prune \) -type f | xargs $svnpath up

		$svnpath revert -R ./Build
		$svnpath up ./Build --username="chenli" --password="chenli"
		$svnpath revert -R ./DepPackages
		$svnpath up ./DepPackages --username="chenli" --password="chenli"
		$svnpath revert -R ./Packages
		$svnpath up ./Packages --username="chenli" --password="chenli"
		$svnpath revert -R ./ProjectSettings
		$svnpath up ./ProjectSettings --username="chenli" --password="chenli"
		$svnpath revert -R ./Tools
		$svnpath up ./Tools --username="chenli" --password="chenli"
		$svnpath revert -R ./UserSettings
		$svnpath up ./UserSettings --username="chenli" --password="chenli"

		#revert and update "Assets", exclude "Assets/Resources"
		cd ./Assets
		echo "cd:"`pwd`
		for file in ./*
		do
			if [[ -d $file ]];
			then
				if [[ $file == "./Resources" ]]
				then
					echo -e "\nnot revert "$file"\n"
				else
					echo "revert -R "$file
					$svnpath revert -R $file"@"
					echo "update "$file
					$svnpath up $file --username="chenli" --password="chenli"
				fi
			else
				if [[ $file == "./Resources.meta" ]]
				then
					echo -e "\nnot revert "$file"\n"
				else
					echo "revert "$file
					$svnpath revert $file"@"
					echo "update "$file
					$svnpath up $file --username="chenli" --password="chenli"
				fi
				
			fi
		done
		#update Resources/.
		$svnpath up ./Resources/Fonts --username="chenli" --password="chenli"
		$svnpath up ./Resources/Shader --username="chenli" --password="chenli"
		$svnpath up ./Resources/UGUI --username="chenli" --password="chenli"
		$svnpath up ./Resources/SDKConfigs --username="chenli" --password="chenli"
		$svnpath up ./Resources/Manager --username="chenli" --password="chenli"
	} #>/dev/null 2>&1
}

#清空xlua目录，解决脚本编译错误问题
delete_xlua_gen()
{
	client_path=$1
	echo -e "\n---------------delete_xlua_gen---------------\n"
	gen_path=$client_path"/Assets/Xlua/Gen"
	if [[ -d "$gen_path" ]]; then
		rm -rf "$gen_path"
		echo "exist xlua path and delete : "$gen_path $?
	fi
	sleep 10
}

svn_up()
{
	client_path=$1
	echo -e "\n...start svn up..."
	cd $client_path
	echo "cd "`pwd`
	$svnpath status
	$svnpath cleanup
	{
		$svnpath revert -R .
		$svnpath up --username="chenli" --password="chenli"
	} #>/dev/null 2>&1
}

override_servertag()
{
	echo -e "\n---------------override_servertag---------------\n"
	file_name=$client_chubao"/Assets/Resources/ServerTag.txt"
	echo "servertag path => "$file_name

	echo "clear $file_name"
	true>$file_name
	echo "ServerTag = "$server_tag
	echo "override $file_name"
	cat>>$file_name<<EOF
$server_tag
EOF
	echo "result = "$?
}

delete_res()
{
	echo -e "\n---------------delete_res---------------\n"
	client_path=$1

	log_file="$log_path""$time_data""_delete_res.log"
	if [[ ! -d "$log_path" ]]; then
		mkdir "$log_path"
	fi
	if [[ -f "$log_file" ]]; then
		rm -rf "$log_file"
	fi

	method_name="Model.BuildGames.DeleteResourceFileResSh"
	echo "method_name:"$method_name
	cmd="$unity_path -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name"
	
	{
		echo "cmd:$cmd"
		"$unity_path" -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name
	}||{ 
		errorlogExit "delete_res failed"
	}
}



build_assets()
{
	echo -e "\n---------------build_assets---------------\n"
	client_path=$1

	log_file="$log_path""$time_data""_build_assets.log"
	if [[ ! -d "$log_path" ]]; then
		mkdir "$log_path"
	fi
	if [[ -f "$log_file" ]]; then
		rm -rf "$log_file"
	fi

	method_name="AssetBuild.AssetBuildTools.BuildAssetBundleSh"
	echo "method_name:"$method_name

	custom_args="buildVersion=$build_version"
	echo "custom_args:"$custom_args
	cmd="$unity_path -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name -CustomArgs $custom_args"
	{
		echo "cmd:$cmd"
		"$unity_path" -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name -CustomArgs $custom_args
	}||{ 
		errorlogExit "build_assets failed"
	}
}

copy_assets()
{
	target_path=$2
	target_res_folder=$target_path"/${platform}"
	if [ -d "${target_res_folder}" ]; then
		echo "${target_res_folder} exist---delete"
		rm -rf "${target_res_folder}"
	fi
	# shellcheck disable=SC2039
	echo -e "\n---------------copy_assets to $target_path---------------\n"
	#资源目录
	project_res_path=${client_chubao%/*}
	echo "project_res_path => "$project_res_path
	assets_path=$1    #"${project_res_path}/Bundle/copy/${res_folder}"
	echo "assets_path => "$assets_path

	if [ ! -d "${assets_path}" ]; then
    mkdir -p "${assets_path}"
  fi

	{
		echo "cp -r ${assets_path} ${target_path}"
		cp -r "${assets_path}" "${target_path}"
	}||{ 
		echo "copy_assets failed"
	}
}


build_xcode()
{
	echo -e "\n---------------build_xcode---------------\n"
	client_path=$1

	log_file="$log_path""$time_data""_build_xcode.log"
	if [ ! -d "$log_path" ]; then
		mkdir "$log_path"
	fi
	if [ -f "$log_file" ]; then
		rm -rf "$log_file"
	fi

	method_name="Model.BuildGames.BuildiOS"
	echo "method_name:"$method_name

	xcode_name=$xcode_name"_"$time_data"_"$build_version"_"$bundle_code"_"$server_tag"_"$identifier
	custom_args="productPath=$output_path#productName=$xcode_name#keystoreName=$keystore_name#keystorePass=$keystore_pass#keyaliasName=$keyalias_name#keyaliasPass=$keyalias_pass#buildOptions=$build_options#versionCode=$build_value#buildVersion=$build_version"
	echo "custom_args:"$custom_args
	cmd="$unity_path -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name -CustomArgs $custom_args"
	{
		echo "cmd:$cmd"
		"$unity_path" -batchmode -nographics -quit -logFile $log_file -projectPath $client_path -executeMethod $method_name -CustomArgs $custom_args
	}||{ 
		errorlogExit "build_apk failed"
	}

	xcode_proj_path=$output_path"/"$xcode_name

	echo "xcode_proj_path => $xcode_proj_path"

  cp /Users/Shared/Documents/pbj "$xcode_proj_path"

  cd "$xcode_proj_path" || exit

  chmod +x pbj && ./pbj
}

exportOptions_path="/Users/packer/Documents/plist/ExportOptions.plist"

build_ipa()
{

	log_file="$log_path""$time_data""buildipa_archive"".log"
	if [ ! -d "$log_path" ]; then
		mkdir "$log_path"
	fi

	printf "\n------------ xcode build ipa: \n%s\n" "$(pwd)"
	#----------------项目自定义部分(自定义好下列参数后再执行该脚本)#
	# (注意: 因为shell定义变量时,=号两边不能留空格,若scheme_name与info_plist_name有空格,脚本运行会失败)
	echo "cd $xcode_proj_path"
	cd "$xcode_proj_path" || exit
	pwd

	# 计时
	SECONDS=0
	# 是否编译工作空间 (例:若是用Cocopods管理的.xcworkspace项目,赋值true;用Xcode默认创建的.xcodeproj,赋值false)
	is_workspace="true"
	# 指定项目的scheme名称
	scheme_name="Unity-iPhone"
	# 工程中Target对应的配置plist文件名称, Xcode默认的配置文件为Info.plist
	# info_plist_name="Info"
	# 指定要打包编译的方式 : Release,Debug
	build_configuration="Release"
	zoo="PhantomBladeEX"

	# ----------------自动打包部分(无特殊情况不用修改)#
	# 导出ipa所需要的plist文件路径 (默认为AdHocExportOptionsPlist.plist)
	ExportOptionsPlistPath="$exportOptions_path"
	# 获取项目名称，我用find方法找到的不对，所以直接写死了
	# project_name=`find . -name *.xcworkspace | awk -F "[/.]" '{print $(NF-1)}'`
	project_name=${scheme_name}
	# 获取版本号,内部版本号,bundleID
	# info_plist_path="$project_name/$info_plist_name.plist"

	# 时间
	DATE=$(date '+%Y-%m-%d-%H-%m-%S')
	# 指定输出ipa路径
	export_path=/Users/Shared/buildIpa/"$scheme_name-$DATE"
	# 指定输出归档文件地址
	export_archive_path="$export_path/$scheme_name.xcarchive"
	# 指定输出ipa地址
	export_ipa_path="$export_path"
	# 指定输出ipa名称
	ipa_name="$zoo"

	# AdHoc,AppStore,Enterprise三种打包方式的区别: http://blog.csdn.net/lwjok2007/article/details/46379945
	echo "------------------------------------------------------"
	echo "\033[32m开始构建项目  \033[0m"
	# 指定输出文件目录不存在则创建
	if [ -d "$export_path" ] ; then
		echo "exist dir:$export_path"
	else
		echo "not exist dir and create dir:$export_path"
		echo "mkdir -pv $export_path"
		mkdir -pv "$export_path"
	fi

	# 编译前清理工程
	xcodebuild clean -project ${project_name}.xcodeproj \
	                 -scheme ${scheme_name} \
	                 -configuration ${build_configuration}

	echo ""

	xcodebuild archive -project ${project_name}.xcodeproj \
	                   -scheme ${scheme_name} \
	                   -configuration ${build_configuration} \
	                   -archivePath "${export_archive_path}" \
	              > "$log_file" 2>&1

	#  检查是否构建成功
	#  xcarchive 实际是一个文件夹不是一个文件所以使用 -d 判断
	if [ -d "$export_archive_path" ] ; then
		echo "\033[32;1m项目构建成功 \033[0m"
	else
		echo "\033[31;1m项目构建失败 \033[0m"
		errorlogExit "build project failed!"
	fi
	echo "------------------------------------------------------"

	log_file1="$log_path""$time_data""buildipa_exportipa"".log"
	echo "\033[32m开始导出ipa文件 \033[0m"
	xcodebuild  -exportArchive \
	            -archivePath "${export_archive_path}" \
	            -exportPath "${export_ipa_path}" \
	            -exportOptionsPlist ${ExportOptionsPlistPath} \
	            -allowProvisioningUpdates > "$log_file1" 2>&1

	# 检查文件是否存在
	if [ -f "$export_ipa_path/$ipa_name.ipa" ] ; then
		echo "\033[32;1m导出 ${ipa_name}.ipa 包成功 \033[0m"
		open "$export_path"
	else
		echo "\033[31;1m导出 ${ipa_name}.ipa 包失败 \033[0m"
		errorlogExit "export ipa failed:${ipa_name}.ipa"
	fi

	# 输出打包总用时
	# shellcheck disable=SC2039
	echo "\033[36;1m使用AutoPackageScript打包总用时: ${SECONDS}s \033[0m"
}

upload_res()
{
	echo -e "\n---------------upload_res---------------\n"

	cd $obs_path

	echo "cd"`pwd`

	log_file="$log_path""$time_data""_uploadres"".log"
	if [[ ! -d "$log_path" ]]; then
		mkdir "$log_path"
	fi
	res_folder=$platform
	cmd="./s3_sync sync --concurrency 10 --exclude-ext .meta ${bundle_res_path}/${res_folder} s3://pbex-cdn/yzr3_ex/res/yzr3HD/${res_folder} cos://pbex-cdn-1308734621/yzr3_ex/res/yzr3HD/${res_folder}"
	echo -e "\n upload cmd : ${cmd}"
	echo "result:"`eval "${cmd}"`$?

	refresh_cache_s3="aws cloudfront create-invalidation --distribution-id E1JM5UBMZPCGGV --paths //yzr3_ex/res/yzr3HD/${res_folder}/*"
	echo -e "\nrefresh tx cache:${refresh_cache_s3}"
	echo "result:"`eval "${refresh_cache_s3}"`

	refresh_cache_tx="tccli cdn PurgePathCache --cli-unfold-argument --Paths http://cdn1.pbex.soulframegame.com/ --FlushType delete"
	echo -e "\nrefresh tx cache:${refresh_cache_tx}"
	echo "result:"`eval "${refresh_cache_tx}"`
	
	echo -e "\n upload end!"
}



build_product()
{
	override_servertag
    #build资源
    build_assets $client_chubao
    #将资源拷贝到上传路径
    res_folder=$platform
    project_res_path=${client_chubao%/*}
    echo "project_res_path => "$project_res_path
    assets_md5_path="${project_res_path}/Bundle/Addressables/${res_folder}"
    copy_assets $assets_md5_path $bundle_res_path

    if [ "${B_build_pkg}" = true ]; then
	    #build xcode
	    build_xcode $client_chubao
    fi

    if [ "${B_Build_Ipa}" = true ]; then
      #build ipa
      build_ipa
    fi

    if [ "${B_upload_res}" = true ]; then
    	upload_res
    fi
}

StartBuild()
{
	#开始
	echo -e "\n### start build xcode ###\n"
	#开始
	echo -e "\n### arrange product settings, please wait... ###"

	if [ ${B_svnup} = true ];then
		echo -e "\n###svn update ,please wait...###"
		svn_up $client_chubao

		#清空xlua
		if [[ ${B_build_bundle} == true ]]; then
			delete_xlua_gen ${client_chubao}
		fi
	fi
	
	echo "Build Unity Project to Android APK, please wait..."
	build_product

	# echo -e "\n### build ipa success ###\n---------"
	
}

StartBuild

