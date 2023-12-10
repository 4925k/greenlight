#!/bin/bash
set -eu

## VARIABLES

# time zone
TIMEZONE=America/New_York

# username
USERNAME=greenlight
OS=$(uname)
ARCH=$(dpkg --print-architecture)

# get database password
read -p "Enter password for greenlight DB user:" DB_PASSWORD

# set locale to english
export LC_ALL=en_US.UTF-8

## SCRIPT


echo "[*] Script starting"

echo "[*] Updating and Upgrading System"

# enable the universe repo
add-apt-repository --yes universe

# update all software packages
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "[*] Setting time and locale"

# set system timezone and install locales
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

echo "[*] Creating user"

# add new user with sudo privileges
useradd --create-home --shell "/bin/bash" --groups sudo ${USERNAME}

# force password reset
passwd --delete ${USERNAME}
chage --lastday 0 ${USERNAME}

echo "[*] Setup ssh and firewall rules"

# copy ssh keys from root to new user
cp /root/.ssh/* /home/${USERNAME}
chown -R ${USERNAME}:${USERNAME} .ssh

# configure firewall to allow SSH, HTTP and HTTPS traffic
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# install fail2ban
apt -y install fail2ban
apt -y install curl

echo "[*] install necessary packages"

# install migrate tools
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/migrate

echo "[*] go-migrate installed"

# install postgres
apt --yes install postgresql

# set up greenlight DB and create user with the password entered earlier
sudo -i -u postgres psql -c "CREATE DATABASE greenlight"
sudo -i -u postgres psql -d greenlight -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d greenlight -c "CREATE ROLE greenlight WITH LOGIN PASSWORD '${DB_PASSWORD}'"
sudo -i -u postgres psql -d greenlight -c "GRANT ALL ON SCHEMA public TO ${USERNAME}"


# add database dsn to environment variable
echo "GREENLIGHT_DB_DSN='postgres://greenlight:${DB_PASSWORD}@localhost/greenlight?sslmode=disable'" >> /etc/environment

echo "[*] postgres setup done"

# install caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt -y install caddy

echo "[*][*] Script complete! Rebooting"
reboot