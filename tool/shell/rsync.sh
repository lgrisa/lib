#-v, --verbose 详细模式输出
#    --info=FLAGS 输出INFO级别
#    --debug=FLAGS 输出DEBUG级别
#    --msgs2stderr 用于调试的特殊输出处理
#-q, --quiet 忽略非error的输出
#--no-motd 忽略Daemon模式的MOTD
#-c, --checksum 让自动跳过基于校验和而非默认的修改时间以及文件大小 使用此选项意味着将禁用快速检查，每个文件都采用计算校验码的方式来检查文件是否一致
#-a, --archive 归档(压缩)模式,表示以递归方式传输文件,并保持所有文件属性等同于-rlptgoD(无 -H,-A,-X)
#--no-OPTION 关闭隐含的选项(例如 --no-D)
#-r, --recursive 对子目录以递归模式处理
#-R, --relative 使用相对路径信息
#--no-implied-dirs 不使用--relative发送隐含的目录
#-b, --backup 创建备份,也就是对于目的已经存在有同样的文件名时,将老的文件重新命名为~filename.可以使用--suffix选项来指定不同的备份文件前缀
#--backup-dir=DIR 将备份文件(如~filename)存放在指定目录下
#--suffix=SUFFIX 定义备份文件前缀,默认是~
#-u, --update 仅仅进行更新,也就是跳过所有已经存在于DST,并且文件时间晚于要备份的文件(不覆盖更新的文件)
#--inplace update destination files in-place (SEE MAN PAGE)
#--append 将数据附加到较短的文件 --append-verify 类似--append,但是对旧数据会计算校验和
#-d, --dirs 不使用递归传输目录
# -l, --links 不处理符号链接(保留符号链接)
#-L, --copy-links 将符号链接处理为具体的文件或者文件夹
#    --copy-unsafe-links 只处理不安全的符号链接
#    --safe-links 忽略不在SRC源目录的符号链接
#    --munge-links munge符号链接使它们更安全(但会无法使用)
#-k, --copy-dirlinks 把指向文件夹的符号链接转换为文件夹
#-K, --keep-dirlinks 把接收端的指向文件夹的符号链接当做文件夹
#-H, --hard-links 保留硬链接 -p, --perms 保留权限
#-E, --executability 保留文件的可执行属性
#--chmod=CHMOD 影响文件或文件夹的属性
#-A, --acls 保留ACLs (代表--perms)
#-X, --xattrs 保留扩展属性
#-o, --owner 保留所有者(仅限superuser)
#-g, --group 保留组
#    --devices 保留设备文件(仅限superuser)
#    --copy-devices 把设备文件内容当做文件一样进行复制处理
#    --specials 保留特殊文件
#-D 和--devices --specials一样
#-t, --times 保留修改时间
#-O, --omit-dir-times 忽略文件夹的修改时间
#-J, --omit-link-times 忽略符号链接的修改时间
#    --super 接收端尝试使用superuser进行操作
#    --fake-super 使用xattrs来存储和恢复权限属性
#-S, --sparse 对稀疏文件进行特殊处理以节省空间
#    --preallocate 在写入前预分配DST文件
#-n, --dry-run 执行一个没有实际更改的试运行,只会显示文件会被如何操作
#-W, --whole-file 拷贝文件，不进行增量检测
#-x, --one-file-system 不要跨越文件系统边界
#-B, --block-size=SIZE 检验算法使用的块尺寸，默认是700字节
#-e, --rsh=COMMAND 指定使用rsh,ssh方式进行数据同步
#    --rsync-path=PROGRAM 指定远程服务器上的rsync命令所在路径
#    --existing 仅仅更新那些已经存在于DST的文件，而不备份那些新创建的文件
#    --ignore-existing 跳过更新已存在于DST的文件
#    --remove-source-files 发送方删除非文件夹的源文件
#    --del --delete-during的一个alias
#    --delete 删除那些DST中SRC没有的文件
#    --delete-before 传输前删除，而非传输过程中
#    --delete-during 在传输过程中删除
#    --delete-delay 在传输过程中确定要删除的,在传输结束后进行删除
#    --delete-after 在传输结束后删除，而非传输过程中
#    --delete-excluded 同样删除接收端那些被该选项指定排除的文件
#    --ignore-missing-args 忽略丢失的源参数不输出错误
#    --delete-missing-args 从DEST删除丢失的源参数
#    --ignore-errors 即使出现I/O错误也进行删除
#    --force 即使文件夹非空也强制删除
#    --max-delete=NUM 不删除超过指定数量的文件
#    --max-size=SIZE 不传输超过指定大小的文件
#    --min-size=SIZE 不传输小于指定大小的文件
#    --partial 保留那些因故没有完全传输的文件,以是加快随后的再次传输(即断点续传)
#    --partial-dir=DIR 将因故没有完全传输的文件放到指定文件夹
#    --delay-updates 在传输末尾把所有更新的文件放到位
#-m, --prune-empty-dirs 从文件列表中删除空目录链
#    --numeric-ids 不要把uid/gid值映射为用户/组名
#    --usermap=STRING 自定义用户名映射
#    --groupmap=STRING 自定义组名映射
#    --chown=USER:GROUP 简单的用户/组名映射
#    --timeout=SECONDS 设置I/O超时,单位为秒
#    --contimeout=SECONDS 设置Daemon连接超时,单位为秒
#-I, --ignore-times 不跳过那些有同样的时间和大小的文件
#-M, --remote-option=OPTION 只把指定选项发送到远端
#    --size-only 只跳过大小相同的文件
#    --modify-window=NUM 决定文件是否时间相同时使用的时间戳窗口，默认为0
#-T, --temp-dir=DIR 在指定文件夹中创建临时文件
#-y, --fuzzy 如果DEST没有任何文件,查找类似的文件
#    --compare-dest=DIR 同样比较DIR中的文件来决定是否需要备份
#    --copy-dest=DIR 和上面的类似,但是还会复制指定文件夹中的没有改变的文件
#    --link-dest=DIR 和上面类似,只是没有改变的文件会被硬链接到DST
#-z, --compress 在传输过程中进行压缩 --compress-level=NUM 指定压缩级别0-9,默认为6
#    --skip-compress=LIST 跳过压缩文件后缀在指定列表中的文件
#-C, --cvs-exclude 自动跳过CVS的生成文件
#-f, --filter=RULE 添加一个文件过滤规则
#-F 等于--filter='dir-merge /.rsync-filter' 重复的: --filter='- .rsync-filter'
#    --exclude=PATTERN 排除符合匹配规则的文件
#    --exclude-from=FILE 从指定文件中读取需要排除的文件
#    --include=PATTERN 包含(不排除)符合匹配规则的文件
#    --include-from=FILE 从指定文件中读取需要包含(不排除)的文件
#    --files-from=FILE 从指定文件中读取SRC源文件列表
#-0, --from0 从文件中读取的文件名以'\0'终止
#-s, --protect-args 没有空格分隔;只有通配符的特殊字符
#    --address=ADDRESS 绑定到指定的地址
#    --port=PORT 指定其他的rsync服务端口
#    --sockopts=OPTIONS 指定自定义的TCP选项
#    --blocking-io 对远程shell使用阻塞IO
#    --stats 提供某些文件的传输状态
#-8, --8-bit-output 在输出中留下高比特的字符
#-h, --human-readable 用人类可读的格式输出数字
#    --progress 在传输过程中显示进度
#-P 等同于--partial --progress
#-i, --itemize-changes 输出对所有更新的变更摘要
#    --out-format=FORMAT 用指定格式输出更新
#    --logutil-file=FILE 将日志保存到指定文件
#    --logutil-file-format=FMT 用指定格式更新日志
#    --password-file=FILE 从文件读取Daemon服务器密码
#    --list-only 不复制而是只列出
#    --bwlimit=RATE 限制套接字I/O带宽
#    --outbuf=N|L|B 设置输出缓冲,为None,Line或者Block
#    --write-batch=FILE 写入批量更新到指定文件
#    --only-write-batch=FILE 和上面类似,但是对DST进行只写的更新
#    --read-batch=FILE 从指定文件读取一个批量更新
#    --protocol=NUM 强制使用指定的老版本协议
#    --iconv=CONVERT_SPEC 对文件名进行字符编码转换
#    --checksum-seed=NUM 设置块/文件的校验和种子
#-4, --ipv4 偏向于使用IPv4 -6,
#    --ipv6 偏向于使用IPv6
#    --version 打印版本号
#(-h) --help 显示帮助信息

start=$(date +%s)

sshpass -p LINGdi1535 rsync -arPz --delete -e "ssh -p 22" --chmod=ugo=rwx Library/ jubin@192.168.0.200:Library/android

end=$(date +%s)

take=$(( end - start ))

echo Time taken to execute commands is ${take} seconds.