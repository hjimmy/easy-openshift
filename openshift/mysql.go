package main

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"net/http"
	"io/ioutil"
	"errors"
	"strings"
	"crypto/tls"
	"gopkg.in/yaml.v2"
	//"os/exec"
        //"os"
)

const KUBE_CONFIG_DEFAULT_LOCATION string = "/etc/origin/master/admin.kubeconfig"

func check(e error) {
    if e != nil {
        panic(e)
    }
}


var Pvc_template = `
  apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    name: owncloud1
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 20Gi
`

type Pvc struct {
     ApiVersion string `yaml:"apiVersion"`
     Kind string `yaml:"kind"`
     Metadata  struct {
          Name string `yaml:"name"`
     }

     Spec struct {
        AccessModes [] string `yaml:"accessModes"`
        Resources struct{
           Requests struct{
                 Storage string
           }
        }
     } `yaml:"spec"`

}



var Imagestream_template = `
  apiVersion: v1
  kind: ImageStream
  metadata:
    labels:
      app: mysql
    name: mysql
  spec:
    tags:
    - from:
        kind: DockerImage
        name: docker.io/mysql:latest
      importPolicy: {}
      name: latest
      referencePolicy:
        type: ""
  status:
    dockerImageRepository: ""
`

type Imagestream struct {
     ApiVersion string `yaml:"apiVersion"`
     Kind    string `yaml:"kind"`
     Metadata  struct {
                Labels struct {
                           App string `yaml:"app"`
                       }
                Name string `yaml:"name"`
     }

     Spec  struct {
		Tags [] struct{
                    From struct { 
                           Kind string `yaml: "kind"`
			   Name string `yaml: "name"`
		    }
      		    ImportPolicy struct {} `yaml: "importPolicy"`
      		    Name string `yaml: "name"`
      		    ReferencePolicy struct {
        		Type string `yaml: "type"`
                    } `yaml: "referencePolicy"`
                }
		
             }

     Status struct {
        DockerImageRepository string  `yaml:"dockerImageRepository"`
     } `yaml:"status"`

}

var Deploymentconfig_template = `
  apiVersion: v1
  kind: DeploymentConfig
  metadata:
    labels:
      app: mysql
    name: mysql
  spec:
    replicas: 2
    selector:
      app: mysql
      deploymentconfig: mysql
    template:
      metadata:
        labels:
          app: mysql
          deploymentconfig: mysql
      spec:
        containers:
        - env:
          - name: MYSQL_PASSWORD
            value: qwer1234
          - name: MYSQL_USER
            value: root
          - name: MYSQL_ROOT_PASSWORD
            value: qwer1234
          image: docker.io/mysql
          name: mysql
          ports:
          - containerPort: 3306
            protocol: TCP
          volumeMounts:
          - name: mysql-persistent-storage
            mountPath: /var/lib/mysql/
        volumes:
        - name: mysql-persistent-storage
          persistentVolumeClaim:
            claimName: nfs-mysql-pvc
    test: false
    triggers:
    - type: ConfigChange
    - imageChangeParams:
        automatic: true
        containerNames:
        - mysql
        from:
          kind: ImageStreamTag
          name: mysql:latest
      type: ImageChange
  status:
    availableReplicas: 0
    latestVersion: 0
    observedGeneration: 0
    replicas: 0
    unavailableReplicas: 0
    updatedReplicas: 0
`

type Deploymentconfig struct {
        ApiVersion string `yaml:"apiVersion"`
        Kind    string `yaml:"kind"`
        Metadata  struct {
                Name string `yaml:"name"`
		Labels struct {
                  App string `yaml:"app"`
                }
        }

        Spec struct {
                     Replicas int32
                     Selector struct {
                          App string `yaml:"app"`
			  Deploymentconfig string `yaml:"deeploymentconfig"`
                     }
                     Template struct {
                          Metadata struct{
                               Labels struct{
                                        App string `yaml:"app"`
					Deploymentconfig string `yaml:"deeploymentconfig"`
                                   } `yaml:"labels"`
                           } `yaml:"metadata"`

                          Spec struct{
                               Containers [] struct{
                                  Image string `yaml:"image"`
				  Name string `yaml:"name"`
                                  Env []struct {
                                    Name string
                                    Value string
                                  } `yaml:"env"`
				  Ports []struct {
                                    ContainerPort int32 `yaml:"containerPort"`
                                    Protocol string  `yaml:"protocol"`
                                  } `yaml:"ports"` 
                                  VolumeMounts []struct {
                                    Name string `yaml:"name"`
                                    MountPath string `yaml:"mountPath"`
                                  } `yaml:"volumeMounts"`
			       }
                               Volumes [] struct {
				  Name string `yaml:"name"`
				  PersistentVolumeClaim struct {
                                         ClaimName string `yaml:"claimName"`
                                  } `yaml:"persistentVolumeClaim"`
                               } `yaml:"volumes"`
		         }
                    }

		    Test bool `yaml:"test"`
                    Triggers []struct{
                        Type string `yaml:"type"`
                        ImageChangeParams struct {
                                Automatic bool `yaml:"automatic"`
                                ContainerNames [] string `yaml:"containerNames"`
                                From struct {
                                        Kind string  `yaml:"kind"`
                                        Name string  `yaml:"name"`
                                } `yaml:"from"`
                        } `yaml:"imageChangeParams"`
                } `yaml:"triggers"`
 
         }

         Status struct {
                AvailableReplicas int32  `yaml:"availableReplicas"`
                LatestVersion int32 `yaml:"latestVersion"`
                ObservedGeneration int32 `yaml:"observedGeneration"`
                Replicas int32 `yaml:"replicas"`
                UnavailableReplicas int32 `yaml:"unavailableReplicas"`
                UpdatedReplicas int32 `yaml:"updatedReplicas"`
           }
}

var Service_template = `
  apiVersion: v1
  kind: Service
  metadata:
    labels:
      app: mysql
    name: mysql
  spec:
    type: NodePort
    ports:
    - name: 3306-tcp
      port: 3306
      protocol: TCP
      nodePort: 30002
    selector:
      app: mysql
      deploymentconfig: mysql
`

type Service struct {
   ApiVersion string `yaml:"apiVersion"`
   Kind    string `yaml:"kind"`
   Metadata  struct {
        Labels struct {
               App string `yaml:"app"`
        } `yaml:"labels"`
        Name string `yaml:"name"`
     } `yaml:"metadata"`
   
    Spec struct {
        Type string `yaml:"type"`
        Ports [] struct{
            Name string `yaml:"name"`
            Port int32  `yaml:"port"`
            Protocol string `yaml:"protocol"`
            NodePort int32 `yaml:"nodePort"`
        }
        Selector struct{
            App string `yaml:"app"`
            Deploymentconfig string `yaml:"deploymentconfig"`
        }
    } 
}

func load_user_token(username string) (ret string, err error){
     
    source, err := ioutil.ReadFile(KUBE_CONFIG_DEFAULT_LOCATION)
    check(err)
    yaml, err := simpleyaml.NewYaml(source)
    check(err)
    size, err := yaml.Get("users").GetArraySize()

    var admin_token string
    for i := 0; i < size; i++ {
        namefull, err := yaml.Get("users").GetIndex(i).Get("name").String()
        check(err)
        name := namefull[:6]
        if name == username + "/" {
           admin_token, err = yaml.Get("users").GetIndex(i).Get("user").Get("token").String()
           return admin_token, nil   
        }                          
     }
     return "", errors.New("User admin token is not exist!")

}

func create_mysql_service(svcname string, serveraddr string, serverport string, pname string){

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/api/v1/namespaces/" +  pname + "/services"

    service := Service{}
    err = yaml.Unmarshal([]byte(Service_template), &service)
    check(err)
    service.Metadata.Labels.App = svcname
    service.Metadata.Name = svcname
    service.Spec.Type = "NodePort"
    service.Spec.Ports[0].NodePort = 30002
    service.Spec.Selector.App = svcname
    service.Spec.Selector.Deploymentconfig = svcname
    service_new, err := yaml.Marshal(&service)
    check(err)
    service_str := string(service_new)
    payload := strings.NewReader(service_str)
    req, _ := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)

    res, _ := client.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))
}


func create_mysql_imagestream(svcname string, serveraddr string, serverport string, pname string){

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)
    url := "https://" + serveraddr + ":" + serverport + "/oapi/v1/namespaces/" +  pname + "/imagestreams"

    imagestream := Imagestream{}
    err = yaml.Unmarshal([]byte(Imagestream_template), &imagestream)
    check(err)
    imagestream.Metadata.Labels.App = svcname
    imagestream.Metadata.Name = svcname

    imagestream_new, err := yaml.Marshal(&imagestream)
    check(err)
    imagestream_str := string(imagestream_new)
    payload := strings.NewReader(imagestream_str)
    req, _ := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)

    res, err := client.Do(req)
    if err != nil {
	    fmt.Println(err)
           return
    }
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(string(body))

}

func create_mysql_deploymentconfig(svcname string, serveraddr string, serverport string, pname string){ 

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/oapi/v1/namespaces/" +  pname + "/deploymentconfigs"

    deploymentconfig := Deploymentconfig{}
    err = yaml.Unmarshal([]byte(Deploymentconfig_template), &deploymentconfig)
    check(err)

    deploymentconfig.Metadata.Name = svcname
    deploymentconfig.Spec.Replicas = 1   

    deploymentconfig.Spec.Template.Metadata.Labels.App = svcname


    deploymentconfig.Spec.Template.Spec.Containers[0].Image = svcname + ":latest"
    deploymentconfig.Spec.Template.Spec.Containers[0].Name = svcname


    deploymentconfig_new, err := yaml.Marshal(&deploymentconfig)
    check(err)
    deploymentconfig_str := string(deploymentconfig_new)
    payload := strings.NewReader(deploymentconfig_str)
    req, _ := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)

    res, _ := client.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(string(body))

}

func create_owncloud_pvc(svcname string, serveraddr string, serverport string, pname string, size string) {

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/api/v1/namespaces/" +  pname + "/persistentvolumeclaims"
    fmt.Println(url)
    pvc := Pvc{}
    err = yaml.Unmarshal([]byte(Pvc_template), &pvc)
    check(err)
    pvc.Spec.Resources.Requests.Storage = size + "Gi"
    pvc.Metadata.Name = "nfs-" + svcname + "-pvc"

    pvc_new, err := yaml.Marshal(&pvc)
    check(err)
    pvc_str := string(pvc_new)
    fmt.Println(pvc_str)
    payload := strings.NewReader(pvc_str)
    req, _ := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)

    res, _ := client.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(string(body))
 }



func main() {
   create_owncloud_pvc("mysql", "10.1.110.161", "8443", "test2", "10")
   create_mysql_imagestream("mysql", "10.1.110.161", "8443", "test2")
   create_mysql_deploymentconfig("mysql", "10.1.110.161", "8443", "test2")
   create_mysql_service("mysql", "10.1.110.161", "8443", "test2")
} 
