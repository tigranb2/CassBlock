cd
sudo apt update
sudo apt -y upgrade
sudo apt install -y build-essential linux-headers-$(uname -r)
sudo apt install -y git zip python3-pip

cd
wget -q https://dl.google.com/go/go1.14.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
rm go1.14.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
echo 'export GOPATH=~/go' >> /root/.bashrc
source /root/.bashrc

cd
wget -q https://github.com/ethereum/go-ethereum/archive/v1.9.12.tar.gz
tar -xzf v1.9.12.tar.gz
rm v1.9.12.tar.gz
cd go-ethereum-1.9.12
make all

cd
echo 'export PATH=$PATH:~/go-ethereum-1.9.12/build/bin' >> /root/.bashrc
source /root/.bashrc

cd
mkdir mininet
cd mininet
git clone git://github.com/mininet/mininet
cd mininet
git checkout -b 2.2.2
cd ..
mininet/util/install.sh -a
sudo mn --test pingall

cd
sudo apt install openjdk-8-jre-headless
echo 'JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64' >> /root/.bashrc
source /root/.bashrc

apt install curl

cd
curl -OL https://downloads.apache.org/cassandra/3.11.10/apache-cassandra-3.11.10-bin.tar.gz
tar -xzvf apache-cassandra-3.11.10-bin.tar.gz
rm apache-cassandra-3.11.10-bin.tar.gz
mv apache-cassandra-3.11.10 cassandra

cd
go get github.com/gocql/gocql
