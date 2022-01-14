export COMPOSER_HOME=/usr

# 用来防止该脚本被重复执行
#check="$0"
#(( procCnt=`ps -A --format='%p%P%C%x%a' --width 2048 -w --sort pid|grep "$check"|grep -v grep|grep -v " -c sh "|grep -v "$$" | grep -c sh|awk '{printf("%d",$1)}'` ))
#if [ ${procCnt} -gt 0 ] ; then
#    echo "$0脚本已经在运行[procs=${procCnt}],此次执行自动取消."
#    exit 1;
#fi

echo "===============$(date +%F_%T) 开始发布===============";

echo -e "\n--- 进入项目目录";
cd "$(pwd)" || exit;
pwd;
echo "---";

echo -e "\n--- 当前分支是";
git branch;
echo '---';

echo -e "\n--- 当前用户是";
whoami;
echo '---';

echo -e "\n--- 更新.env文件";
cp .env.example .env;
while [ -n "$1" ]
do
  case "$1" in
    -APP_URL)
        echo "发现 -APP_URL 选项，值是：$2"
        sed -i "s/{{APP_URL}}/$2/" .env
        shift
        ;;
    -DB_DATABASE)
        echo "发现 -DB_DATABASE 选项，值是：$2"
        sed -i "s/{{DB_DATABASE}}/$2/" .env
        shift
        ;;
    -DB_USERNAME)
        echo "发现 -DB_USERNAME 选项，值是：$2"
        sed -i "s/{{DB_USERNAME}}/$2/" .env
        shift
        ;;
    -REDIS_PASSWORD)
        echo "发现 -REDIS_PASSWORD 选项，值是：$2"
        sed -i "s/{{REDIS_PASSWORD}}/$2/" .env
        shift
        ;;
    -DB_PASSWORD)
        echo "发现 -DB_PASSWORD 选项，值是：$2"
        sed -i "s/{{DB_PASSWORD}}/$2/" .env
        shift
        ;;
    -BEIANCODE)
        echo "发现 -BEIANCODE 选项，值是：$2"
        sed -i "s/{{BEIANCODE}}/$2/" .env
        shift
        ;;
    *)
        echo "不支持这个选项：$1"
        ;;
  esac
  shift
done

echo '---'

echo -e "\n--- 最近的提交信息->$(git log --pretty=format:"%an:%s" -1) ---";

echo -e "\n--- 安装composer依赖";
composer install --ignore-platform-reqs --optimize-autoloader --no-dev 2>&1;
echo '---';

echo -e "\n--- 清除缓存并执行数据库迁移";
php artisan optimize:clear && php artisan migrate --force 2>&1;
echo '---';

echo -e "\n--- 安装npm依赖";
npm install
echo '---';

echo -e "\n--- 构建静态资源";
npm run production
echo '---';

echo -e "---- 当前目录：$(pwd) ----\r\n"

echo -e "---- 清理路由缓存后查看opcache状态"
php artisan route:clear && php artisan opcache:status
echo -e "----\r\n"

# 清理opcache
echo -e "---- 清理路由缓存和opcache缓存"
php artisan route:clear && php artisan opcache:clear
echo -e "----\r\n"

# 生成缓存
echo -e "---- 生成项目缓存"
php artisan config:cache
php artisan route:cache
php artisan view:cache
echo -e "----\r\n"

echo -e "---- opcache预编译"
php artisan opcache:compile --force
echo -e "----\r\n"

echo -e "\n--- 清理页面静态化缓存"
php artisan page-cache:clear
echo -e "---\n";

echo -e "\n--- 停止Horizon（停止后应该由supervisor再次启动）"
php artisan horizon:terminate
echo -e "---\n";

echo -e "\n--- 改文件夹权限"
chown -R www:www ./
echo -e "---\n";

echo "===============$(date +%F_%T) 结束发布===============";
