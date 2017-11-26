

#启动redis服务器
redis-server ./conf/redis.conf &

#启动fdfs_trackerd
/etc/fdfs/fdfs_trackerd ./conf/tracker.conf restart
#启动fdfs_stroaged
/etc/fdfs/fdfs_storaged ./conf/storage.conf restart

