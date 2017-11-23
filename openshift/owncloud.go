package openshift

import (
	"fmt"
	"github.com/smallfish/simpleyaml"
	"net/http"
	"io/ioutil"
	"errors"
	"strings"
	"crypto/tls"
	"gopkg.in/yaml.v2"
        "strconv"
//	"os/exec"
        //"os"
)

const KUBE_CONFIG_DEFAULT_LOCATION string = "/etc/origin/master/admin.kubeconfig"
const serveraddr string = "127.0.0.1"
const serverport string = "8443"


func check(e error) {
    if e != nil {
        panic(e)
    }
}

var Project_template = `
  apiVersion: v1
  kind: Project
  metadata:
    name: test3 
`

type Project struct {
     ApiVersion string `yaml:"apiVersion"`
     Kind    string `yaml:"kind"`
     Metadata  struct {
         Name string `yaml:"name"`
     } `yaml:"metadata"`
}

var Imagestream_template = `
  apiVersion: v1
  kind: ImageStream
  metadata:
    labels:
      app: owncloud
    name: owncloud
  spec:
    tags:
    - from:
        kind: DockerImage
        name: docker.io/owncloud:latest
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
      app: owncloud
    name: owncloud
  spec:
    replicas: 2
    selector:
      app: owncloud
      deploymentconfig: owncloud
    template:
      metadata:
        labels:
          app: owncloud
          deploymentconfig: owncloud
      spec:
        containers:
        - image: owncloud
          name: owncloud
          ports:
          - containerPort: 80
            protocol: TCP

          volumeMounts:
          - name: owncloud
            mountPath: /var/www/html/data/
        volumes:
        - name: owncloud
          persistentVolumeClaim:
              claimName: nfs-owncloud-pvc

      restartPolicy: Always

    test: false
    triggers:
    - type: ConfigChange
    - imageChangeParams:
        automatic: true
        containerNames: []
        from:
          kind: ImageStreamTag
          name: owncloud:latest
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
		Labels struct {
    		           App string `yaml:"app"`
		       }
                Name string `yaml:"name"`
        }
	
	Spec struct {
		     Replicas int
		     Selector struct {
			    App string `yaml:"app"`
                            Deploymentconfig string `yaml:"deploymentconfig"`
		     }
    		     Template struct {
                          Metadata struct{
        		       Labels struct{
          				App string `yaml:"app"`
          		                Deploymentconfig string `yaml:"deploymentconfig"`
                                   }
			   }

                          Spec struct{
		               Containers [] struct{
			       	  Image string `yaml:"image"`
			          Name string  `yaml:"name"`
		                  Ports []struct {
			            ContainerPort int `yaml:"containerPort"`
            		            Protocol string `yaml:"protocol"`
			  	  }`yaml:"ports"`
                                 VolumeMounts []struct {
                                    Name string `yaml:"name"`
                                    MountPath string `yaml:"mountPath"`
                                  } `yaml:"volumeMounts"`
                               }
                               Volumes [] struct{
				   Name string `yaml:"name"`
                                   PersistentVolumeClaim struct {
                                          ClaimName string `yaml:"claimName"`
				   }`yaml:"persistentVolumeClaim"`
		
                                 } `yaml:"volumes"`
                               
                           }
	
                         RestartPolicy string `yaml:"restartPolicy"`
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
    		AvailableReplicas int  `yaml:"availableReplicas"`
    		LatestVersion int `yaml:"latestVersion"`
    		ObservedGeneration int `yaml:"observedGeneration"`
    		Replicas int `yaml:"replicas"`
    		UnavailableReplicas int `yaml:"unavailableReplicas"`
    		UpdatedReplicas int `yaml:"updatedReplicas"`
	   }
  }


var Service_template =`
  apiVersion: v1
  kind: Service
  metadata:
    labels:
      app: owncloud
    name: owncloud
  spec:
    ports:
    - name: 80-tcp
      port: 80
      protocol: TCP
      nodePort: 1000
    selector:
      app: owncloud
      deploymentconfig: owncloud
`

type Service struct{
     ApiVersion string `yaml:"apiVersion"`
     Kind    string `yaml:"kind"`
     Metadata  struct {
                Labels struct {
                           App string `yaml:"app"`
                       }
                Name string `yaml:"name"`
     }
     Spec struct {
        Type string `yaml:"type"`
	Ports [] struct{
	    Name string `yaml:"name"`
	    Port int  `yaml:"port"`
	    Protocol string `yaml:"protocol"`
            NodePort int `yaml:"nodePort"`
	}
	Selector struct{
	    App string `yaml:"app"`
	    Deploymentconfig string `yaml:"deploymentconfig"`	
        }
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

func create_owncloud_service(svcname string, pname string, nodeport int){

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
    service.Spec.Ports[0].NodePort = nodeport
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


func create_owncloud_imagestream(svcname string, pname string){

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
    fmt.Println(imagestream_str)
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

func create_owncloud_deploymentconfig(svcname string, pname string, replica int){ 

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/oapi/v1/namespaces/" +  pname + "/deploymentconfigs"

    deploymentconfig := Deploymentconfig{}
    err = yaml.Unmarshal([]byte(Deploymentconfig_template), &deploymentconfig)
    fmt.Println(err)
    check(err)
    deploymentconfig.Metadata.Labels.App = svcname
    deploymentconfig.Metadata.Name = svcname
    deploymentconfig.Spec.Selector.App = svcname
    deploymentconfig.Spec.Replicas = replica
    deploymentconfig.Spec.Selector.Deploymentconfig = svcname
    deploymentconfig.Spec.Template.Metadata.Labels.App = svcname
    deploymentconfig.Spec.Template.Metadata.Labels.Deploymentconfig = svcname
    deploymentconfig.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName = svcname
    deploymentconfig.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = svcname

    deploymentconfig.Spec.Template.Spec.Containers[0].Image = svcname + ":latest"
    deploymentconfig.Spec.Template.Spec.Containers[0].Name = svcname

    deploymentconfig.Spec.Triggers[1].ImageChangeParams.ContainerNames = append(deploymentconfig.Spec.Triggers[1].ImageChangeParams.ContainerNames, svcname)
    deploymentconfig.Spec.Triggers[1].ImageChangeParams.From.Name = svcname + ":latest"

    deploymentconfig_new, err := yaml.Marshal(&deploymentconfig)
    check(err)
    deploymentconfig_str := string(deploymentconfig_new)
    fmt.Println(deploymentconfig_str)
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

func create_owncloud_pvc(svcname string, pname string, size int) {

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/api/v1/namespaces/" +  pname + "/persistentvolumeclaims"
    fmt.Println(url)
    pvc := Pvc{}
    err = yaml.Unmarshal([]byte(Pvc_template), &pvc)
    check(err)
    pvc.Spec.Resources.Requests.Storage = strconv.Itoa(size) + "Gi"
    pvc.Metadata.Name = svcname

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

func create_owncloud_project(pname string) {

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := load_user_token("admin")
    check(err)

    url := "https://" + serveraddr + ":" + serverport + "/oapi/v1/projects"
    fmt.Println(url)
    project := Project{}
    err = yaml.Unmarshal([]byte(Project_template), &project)
    check(err)
    project.Metadata.Name = pname

    project_new, err := yaml.Marshal(&project)
    check(err)
    project_str := string(project_new)
    fmt.Println(project_str)
    payload := strings.NewReader(project_str)
    req, _ := http.NewRequest("POST", url, payload)

    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)

    res, _ := client.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    fmt.Println(string(body))

 }


func create_owncloud(appname string, nodeport int, size int, pname string, replica int) {

   create_owncloud_pvc(appname, pname, size)
   create_owncloud_imagestream(appname, pname)
   create_owncloud_deploymentconfig(appname, pname, replica)
   create_owncloud_service(appname, pname, nodeport)
} 


func main(){
    var nodeport int
    var size int
    var replica int

    appname := "owncloud"
    nodeport = 30001
    size = 10
    pname := "owncloud"
    replica = 2
 
    create_owncloud(appname, nodeport, size, pname, replica)

}
