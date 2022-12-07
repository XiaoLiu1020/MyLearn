- [`nginx.conf`](#nginxconf)
- [下面是工作中的使用过的`demo`](#下面是工作中的使用过的demo)
- [`proxy_mall`](#proxy_mall)
- [`production_proxy_mall_conf`](#production_proxy_mall_conf)
- [`production_proxy_pay.conf`](#production_proxy_payconf)
- [`production_page_mall_admin.conf`](#production_page_mall_adminconf)
- [`production-socket.conf`](#production-socketconf)

# `nginx.conf`

```nginx
# For more information on configuration, see:
#   * Official English Documentation: https://nginx.org/en/docs/
#   * Official Russian Documentation: https://nginx.org/ru/docs/

user root;
worker_processes 8;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

events {
    worker_connections 10240;
}

http {
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    include             /etc/nginx/mime.types;
    default_type        application/octet-stream;
    index               index.html;

    gzip on;
    gzip_min_length 1024;
    #gzip_buffers 4 16k;
    #gzip_comp_level 5;
    #gzip_types text/plain application/x-javascript text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;
    gzip_types      text/plain text/css text/csv application/javascript application/json;
    #gzip_disable "MSIE [1-6]\.";

    # Load modular configuration files from the /etc/nginx/conf.d directory.
    # See https://nginx.org/en/docs/ngx_core_module.html#include
    # for more information.
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*.conf;
    include /etc/nginx/http.d/*.conf;
    include /etc/nginx/https.d/*.conf;

    deny 120.78.51.47;
}

stream {
        include /etc/nginx/tcp.d/*.conf;
}

```
# 下面是工作中的使用过的`demo`

`single_page_web_fontend`

```bash
server {
        listen      80;
        server_name tesing.example.com;

        root /home/test/pages/mall_admin_phone/current;

        location / {
                try_files $uri $uri/ /index.html;
        }

        location ~ MP_verify_(.*)\.txt$ {		## $1表示前面（）之内的内容
                add_header Content-Type "text/plain;charset=utf-8";
                return 200 $1;
        }
```

# `proxy_mall`

```bash
server {
	listen      80;
	server_name testing.example.com;


	location / {
		proxy_pass       http://127.0.0.1:8001;
		proxy_redirect   off;
		proxy_set_header Host              $host;
		proxy_set_header X-Real-IP         $remote_addr;
		proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
	}

	location ~ MP_verify_(.*)\.txt$ {		## $1表示前面（）之内的内容
			add_header Content-Type "text/plain;charset=utf-8";
			return 200 $1;
	}
```

# `production_proxy_mall_conf`

```nginx
server {
        listen 443 ssl;
        server_name api.example.com;

        ssl_certificate     /etc/nginx/ssl/202302/7316403_api.example.com.pem;
        ssl_certificate_key /etc/nginx/ssl/202302/7316403_api.example.com.key;

        root /home/webapi/v2/current/app;

        location /v2 {
                proxy_pass http://127.0.0.1:5020/v2;
                proxy_redirect off;
                proxy_set_header Host                   $host;
                proxy_set_header X-Real-IP              $remote_addr;
                proxy_set_header X-Forwarded-For        $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto      $scheme;
        }

        location /static {
                root /home/webapi/v2/current/app;
        }

        location /bridge {
                proxy_pass http://127.0.0.1:5900/bridge;
                proxy_redirect off;
                proxy_set_header Host                   $host;
                proxy_set_header X-Real-IP              $remote_addr;
                proxy_set_header X-Forwarded-For        $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto      $scheme;
        }

        location /app_version {
                proxy_pass http://127.0.0.1:9000;
                proxy_redirect off;
                proxy_set_header Host                   $host;
                proxy_set_header X-Real-IP              $remote_addr;
                proxy_set_header X-Forwarded-For        $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto      $scheme;
        }
        location /download {
                root /home/webapi/data/exports/accounts;
        }

}

```

# `production_proxy_pay.conf`

```nginx
server {
        listen 443 ssl;
        server_name p.example.com;

        root /home/home/pay/v1/current;

        ssl_certificate     /etc/nginx/ssl/202302/7316418_p.example.com.pem;
        ssl_certificate_key /etc/nginx/ssl/202302/7316418_p.example.com.key;

        location / {
                proxy_pass http://127.0.0.1:5002;
                proxy_redirect off;
                proxy_set_header Host                   $host;
                proxy_set_header X-Real-IP              $remote_addr;
                proxy_set_header X-Forwarded-For        $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto      $scheme;
        }

        location /v1/ums/notify {
                proxy_pass http://127.0.0.1:5004/ums/notify;
                proxy_redirect off;
                proxy_set_header Host                   $host;
                proxy_set_header X-Real-IP              $remote_addr;
                proxy_set_header X-Forwarded-For        $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto      $scheme;
        }

        location /static {
                root /home/home/pay/v1/current/app;
        }
}

```

# `production_page_mall_admin.conf`

```nginx
server {
        listen 443 ssl default_server;
        server_name mall.example.com;

        ssl_certificate     /etc/nginx/ssl/202302/7316434_mall.example.com.pem;
        ssl_certificate_key /etc/nginx/ssl/202302/7316434_mall.example.com.key;


        root  /home/app/pages/mall_admin/current;

        location / {
                expires 1d;
                try_files $uri $uri/ /index.html;
        }

        location ~* \.(?:css|js|jpg|jpeg|png|mp4)$ {
                expires 1M;
        }
}
```

# `production-socket.conf`

```nginx
upstream socket {
        hash   $remote_addr consistent;
        server 127.0.0.1:20443;
}

server {
        listen     10443;
        proxy_pass socket;
```

