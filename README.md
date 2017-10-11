Easy OpenShift 
==============================


This will provide a simple way to deploy and use openshift docker container. This project is based on openshift api and beego framwork which can be more easy to ordinary users.

INSTALL:

  os:

     This project should run on Centos 7.2-1611

     You need to install go and mysql/mariadb


  1) Get code:

        $ go get github.com/hjimmy/easy-openshift

	$ cd easy-openshift

	$ cd  conf/app.conf

         Modify  your  config

  2) Init database:

	$ mysql -u $username -p$password

        $ create database easy

        $ mysql -u $username -p$password -D easy < install.sql

  3) Run:

	$ go build

        $ ./run.sh start/stop
	
       
  4) Visit:

      http://localhost:8080

      account：admin   password：123456                                       
