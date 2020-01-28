## Service template

# Information
* GoLand
* Go 1.12.4+ and Go Modules
* Service run on port 3000 (Port 4040 on docker container)
* Create config/main.yaml base on main-template.yaml before you run this template
* Remove unused files and services based on your needs

# Go configs (For Linux)
````$xslt
export GOROOT=/usr/local/go
export GOPATH=~/go
export GOBIN=${GOPATH}/bin
export PATH=${PATH}:/usr/local/bin:${GOROOT}/bin:${GOBIN}
````

# Google credential
* Login Google Develop
* Download credential (JSON)
* Setup ENV in GoLand
* Or with PowerShell:
````$xslt
$env:GOOGLE_APPLICATION_CREDENTIALS="C:\Users\nbxtr\Documents\project\template-golang.json"
````

# Debian packages
````$xslt
sudo apt-get install gcc resolvconf docker.io docker-compose -y
````

# Anaconda environment
````$xslt
conda create -yn backend-golang go
conda activate backend-golang
conda install gxx_linux-64
````

# Run develop mode
````$xslt
GO111MODULE=on go run main.go
````

# Build and run production mode
````$xslt
GO111MODULE=on go build
./backend-golang
````

# Terminate service
````$xslt
sudo kill -9 $(sudo lsof -t -i:3000)
````

# OpenSSL Generator
* In "key" folder
````$xslt
openssl genrsa -out app.rsa 4096
openssl rsa -in app.rsa -pubout > app.rsa.pub
````

# Deploy to Heroku
* 'heroku: true' in main.yaml
````$xslt
heroku login
heroku git:remote --app {HEROKU_APP_NAME}
git push heroku master
heroku logs --tail --app {HEROKU_APP_NAME}
````

# Install Kafka
````$xslt
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure --prefix /usr
make
sudo make install
````

# Generate a self-signed X.509 TLS certificate
Run the following command to generate cert.pem and key.pem files in key folder:
````$xslt
cd key
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
````

# Docker build
````$xslt
docker build --rm -t backend-golang .
````

# Docker run
````$xslt
docker run -d -p 4040:3000 --name backend-golang backend-golang:latest
````

# System troubleshooting
1/ When you has error "cannot unmarshal DNS message":

* Install the resolvconf package.
````$xslt
sudo apt-get purge resolvconf -y
sudo apt-get install resolvconf -y
````

* Edit /etc/resolvconf/resolv.conf.d/head and add the following:
````$xslt
nameserver 8.8.4.4
nameserver 8.8.8.8
````

* Restart the resolvconf service.
````$xslt
sudo service resolvconf restart
````