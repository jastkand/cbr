# CBR Currency Exchange Rate

The original API uses XML format. This service proxies requests to original API and returns JSON.

## Deploy

### Install `supervisor`

```bash
sudo apt update && sudo apt install supervisor
```

```conf
; /etc/supervisor/conf.d/cbr.conf
[program:cbr]
directory=/var/www/app
command=/var/www/app/cbr
autostart=true
autorestart=true
stderr_logfile=/var/log/cbr.err.log
stdout_logfile=/var/log/cbr.out.log
```

```bash
sudo supervisorctl reread
sudo supervisorctl update
```

### Add after deploy script to restart services

```sh
# /var/www/app/after_deploy.sh
sudo supervisorctl stop cbr
mv /var/www/app/cbr-linux /var/www/app/cbr
sudo supervisorctl start cbr
```

```bash
chmod +x /var/www/app/after_deploy.sh
```

### Setup deploy environment variables

Copy `.env.deploy.example` to `.env.deploy` and fill in the values

### Do deploy

From local machine

```bash
make deploy
```
