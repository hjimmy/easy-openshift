//package openshift
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
        "strconv"
	"container/list"
//	"os/exec"
        //"os"
)


type Update_svc_callback func (appname string, pname string, objname string, port int) string
type Update_replica_callback func (appname string, pname string, objname string, replica int)string
type Init_deploymentconfig_callback func (appname string, pname string, replica int) string
type Init_service_callback func (appname string, pname string, nodeport int) string
type Init_pvc_callback func (appname string, pname string, size int) string
type Init_imagestream_callback func (appname string, pname string) string


const KUBE_CONFIG_DEFAULT_LOCATION string = "/etc/origin/master/admin.kubeconfig"
const Serveraddr string = "127.0.0.1"
const Serverport string = "8443"

func Check(e error) {
    if e != nil {
        panic(e)
    }
}

var Scale_template = `
apiVersion: extensions/v1beta1
kind: Scale
metadata:
  name: owncloud
  namespace: owncloud
spec:
  replicas: 4 
`    

type Scale struct {
     ApiVersion string `yaml:"apiVersion"`
     Kind    string `yaml:"kind"`
     Metadata struct {
         Name string `yaml:"name"`
	 Namespace string `yaml:"namespace"`
     } `yaml:"metadata"`
     Spec  struct {
         Replicas int `yaml:"replicas"`
     } `yaml:"spec"`
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
/* 加入test1， test2， test3 是为了兼容mysql， template 不允许有tab符号，否则出错*/
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

          env:
          - name: test1
            value: test
          - name: test2
            value: test
          - name: test3
            value: test

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
                                  Env [] struct{
				    Name string `yaml:"name"`
				    Value string `yaml:"value"`
				  }`yaml:"env"`
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
    resourceVersion: '1110605'
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

type ServiceUpdate struct{
     ApiVersion string `yaml:"apiVersion"`
     Kind    string `yaml:"kind"`
     Metadata  struct {
                Labels struct {
                           App string `yaml:"app"`
                       }   
                Name string `yaml:"name"`
		ResourceVersion string `yaml:"resourceVersion"`
		
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
        ClusterIP string `yaml:"clusterIP"`
      } 
}


var Pvc_template = `
  apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    name: owncloud1
    namespace: owncloud
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
	  Namespace string `yaml:"namespace"`
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
// 创建时因为让它自行选择存储，所以不填写volumename

type PvcUpdate struct {
     ApiVersion string `yaml:"apiVersion"`
     Kind string `yaml:"kind"`
     Metadata  struct {
          Name string `yaml:"name"` 
          Namespace string `yaml:"namespace"`
     }

     Spec struct {
        AccessModes [] string `yaml:"accessModes"`
        Resources struct{
           Requests struct{
                 Storage string
           }
        }
        VolumeName string `yaml:"volumeName"`  // 比pvc多了一个VolumeName，修改容量大小时需要这个值。volumename对应pv的名称。
     } `yaml:"spec"`
      
} 


func Load_user_token(username string) (ret string, err error){
     
    source, err := ioutil.ReadFile(KUBE_CONFIG_DEFAULT_LOCATION)
    Check(err)
    yaml, err := simpleyaml.NewYaml(source)
    Check(err)
    size, err := yaml.Get("users").GetArraySize()

    var admin_token string
    for i := 0; i < size; i++ {
        namefull, err := yaml.Get("users").GetIndex(i).Get("name").String()
        Check(err)
        name := namefull[:6]
        if name == username + "/" {
           admin_token, err = yaml.Get("users").GetIndex(i).Get("user").Get("token").String()
           return admin_token, nil   
        }                          
     }
     return "", errors.New("User admin token is not exist!")

}


func Get_rcname(appname string, pname string) string{
    body := Get_obj(appname, pname, appname, "replicationcontrollers")
    yaml, err := simpleyaml.NewYaml(body)
    Check(err)
    list1, err := yaml.Get("items").Array()
    Check(err)
    len := len(list1)
    for i := 0; i < len; i++ {
       name, err:= yaml.Get("items").GetIndex(i).Get("metadata").Get("name").String()
       Check(err)
       if strings.HasPrefix(name, appname){
             return name	
       }
    }
    return ""
}

func Get_podlist(appname string, pname string, objname string, rcname string) (podlist *list.List){

    body := Get_obj(appname, pname, objname, "pods")
    yaml, _ := simpleyaml.NewYaml(body)
    podarray, err := yaml.Get("items").Array()
    Check(err)
    podlen := len(podarray)
    podlist = list.New()
    for i := 0; i < podlen; i++ {
       pod_name, err:= yaml.Get("items").GetIndex(i).Get("metadata").Get("name").String()
       Check(err)
       if strings.HasPrefix(pod_name, rcname){
             podlist.PushBack(pod_name)
       }

    }
    return podlist
}


/* Update service must kown service's latest version, so we need to get service to get latest version */
func Get_service_rversion(appname string, pname string, objname string) (service ServiceUpdate){
    body := Get_obj(appname, pname, objname, "services") 
    err := yaml.Unmarshal([]byte(body), &service)
    Check(err) 
    return service
}

/*
objtype:
   pods
   deploymentconfigs
   imagestreams
   services
   replicationcontrollers
   persistentvolumeclaims

*/

func Get_obj(appname string, pname string, objname string, objtype string) []byte{
    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := Load_user_token("admin")
    Check(err)
    var url string
    switch objtype{
        case "deploymentconfigs":
            url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
        case "imagestreams":
            url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
	case "services":
	    url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/services/" + objname           
	case "replicationcontrollers":
	    url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype
	case "pods":
	    url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype                
	case "persistentvolumeclaims":
	    url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype + "/" + objname
        default:
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
    }
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)
    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))
    return body
}



func Delete_obj(appname string, pname string, objname string, objtype string){
    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := Load_user_token("admin")
    Check(err)
    var url string
    switch objtype{
        case "deploymentconfigs":
	  url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
        case "imagestreams":
          url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
	default:
          url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname 
    }
    req, _ := http.NewRequest("DELETE", url, nil)
    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)
    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))

}


func Update_svc(appname string, pname string, objname string, port int) string{

    service := Get_service_rversion(appname, pname, objname)
    service.Spec.Ports[0].NodePort = port

    service_new, err := yaml.Marshal(&service)
    Check(err)
    service_str := string(service_new)
    return service_str
}  

func Update_replica(appname string, pname string, objname string, replica int)string{
    scale := Scale{}
    err := yaml.Unmarshal([]byte(Scale_template), &scale)
    Check(err)
    scale.Metadata.Name = appname
    scale.Metadata.Namespace = pname
    scale.Spec.Replicas = replica

    scale_new, err := yaml.Marshal(&scale)
    Check(err)
    scale_str := string(scale_new)
    return scale_str
}


func Init_deploymentconfig(appname string, pname string, replica int) string{

    deploymentconfig := Deploymentconfig{}
    err := yaml.Unmarshal([]byte(Deploymentconfig_template), &deploymentconfig)
    Check(err)
    deploymentconfig.Metadata.Labels.App = appname
    deploymentconfig.Metadata.Name = appname
    deploymentconfig.Spec.Selector.App = appname
    deploymentconfig.Spec.Replicas = replica
    deploymentconfig.Spec.Selector.Deploymentconfig = appname
    deploymentconfig.Spec.Template.Metadata.Labels.App = appname
    deploymentconfig.Spec.Template.Metadata.Labels.Deploymentconfig = appname
    deploymentconfig.Spec.Template.Spec.Volumes[0].Name = appname
    deploymentconfig.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName = appname
    deploymentconfig.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = appname
    switch pname{
	case "owncloud":
            deploymentconfig.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = "/var/www/html/data"
    	case  "mysql":
            deploymentconfig.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath = "/var/lib/mysql"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[0].Name = "MYSQL_PASSWORD"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[0].Value = "qwer1234"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[1].Name = "MYSQL_ROOT_PASSWORD"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[1].Value = "qwer1234"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[2].Name = "MYSQL_USER"
            deploymentconfig.Spec.Template.Spec.Containers[0].Env[2].Value = "root"
            deploymentconfig.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 3306
    }
    deploymentconfig.Spec.Template.Spec.Containers[0].Image = appname + ":latest"
    deploymentconfig.Spec.Template.Spec.Containers[0].Name = appname

    deploymentconfig.Spec.Triggers[1].ImageChangeParams.ContainerNames = append(deploymentconfig.Spec.Triggers[1].ImageChangeParams.ContainerNames, appname)
    deploymentconfig.Spec.Triggers[1].ImageChangeParams.From.Name = appname + ":latest"

    deploymentconfig_new, err := yaml.Marshal(&deploymentconfig)
    Check(err)
    deploymentconfig_str := string(deploymentconfig_new)
    return deploymentconfig_str
}

func Init_pvc(appname string, pname string, size int) string{
    pvc := Pvc{}
    err := yaml.Unmarshal([]byte(Pvc_template), &pvc)
    Check(err)
    pvc.Spec.Resources.Requests.Storage = strconv.Itoa(size) + "Gi"
    pvc.Metadata.Name = appname
    pvc.Metadata.Namespace = pname

    pvc_new, err := yaml.Marshal(&pvc)
    Check(err)
    pvc_str := string(pvc_new)
    return pvc_str
}

func Init_imagestream(appname string, pname string) string{

    imagestream := Imagestream{}
    err := yaml.Unmarshal([]byte(Imagestream_template), &imagestream)
    Check(err)
    switch pname {
    case "owncloud":
        imagestream.Spec.Tags[0].From.Name = "docker.io/owncloud:latest"
    case "mysql":
        imagestream.Spec.Tags[0].From.Name = "docker.io/mysql:latest"
    }
    imagestream.Metadata.Labels.App = appname
    imagestream.Metadata.Name = appname

    imagestream_new, err := yaml.Marshal(&imagestream)
    Check(err)
    imagestream_str := string(imagestream_new)
    return imagestream_str

}


func Init_service(appname string, pname string, nodeport int) string{
    service := Service{}
    err := yaml.Unmarshal([]byte(Service_template), &service)
    Check(err)
    service.Metadata.Labels.App = appname
    service.Metadata.Name = appname
    service.Spec.Type = "NodePort"
    service.Spec.Ports[0].NodePort = nodeport
    switch pname{
	case "owncloud": 
	    service.Spec.Ports[0].Name = "80-tcp"
            service.Spec.Ports[0].Port = 80
        case  "mysql":
            service.Spec.Ports[0].Name = "3306-tcp"
            service.Spec.Ports[0].Port = 3306
    }
    service.Spec.Selector.App = appname
    service.Spec.Selector.Deploymentconfig = appname
    service_new, err := yaml.Marshal(&service)
    Check(err)
    service_str := string(service_new)
    return service_str

}



/*
 objtype: 
    replica: Change pod replicas
    port: Change visit port
*/

func Update_obj(appname string, pname string, replica int, port int, objname string, objtype string){

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := Load_user_token("admin")
    Check(err)
    var url string
    var update string
    switch objtype {
        case "replica":
          url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname + "/deploymentconfigs/" + objname + "/scale"
          var callback Update_replica_callback = Update_replica
          update = callback(appname, pname, objname,replica)
        case "port":
          url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/services/" + appname
	  var callback Update_svc_callback = Update_svc
          update = callback(appname, pname, objname, port)
        default:
          url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname +  "/" + objtype  + "/" + objname
    }

    payload := strings.NewReader(update) 
    req, _ := http.NewRequest("PUT", url, payload)
    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)
    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))

}

/*
objtype:
   pods
   deploymentconfigs
   imagestreams
   services
   replicationcontrollers
   persistentvolumeclaims

*/

func Create_obj(appname string, pname string, replica int, port int, size int, objtype string){

    tr := &http.Transport{ TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
    client := &http.Client{Transport: tr}
    admin_token,err := Load_user_token("admin")
    Check(err)
    var url string
    var create_str string
    switch objtype {
        case "deploymentconfigs":
            url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype
            var callback Init_deploymentconfig_callback = Init_deploymentconfig
            create_str = callback(appname, pname, replica)
        case "imagestreams":
            url = "https://" + Serveraddr + ":" + Serverport + "/oapi/v1/namespaces/" +  pname +  "/" + objtype
	    var callback Init_imagestream_callback = Init_imagestream
            create_str = callback(appname, pname)
        case "services":
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype 
	    var callback Init_service_callback = Init_service
            create_str = callback(appname, pname, port)
        case "replicationcontrollers":
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype
        case "pods":
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype
        case "persistentvolumeclaims":
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname + "/" + objtype
	    var callback Init_pvc_callback = Init_pvc
            create_str = callback(appname, pname, size)
        default:
            url = "https://" + Serveraddr + ":" + Serverport + "/api/v1/namespaces/" +  pname +  "/" + objtype
    }

    payload := strings.NewReader(create_str)
    req, _ := http.NewRequest("POST", url, payload)
    req.Header.Add("content-type", "application/yaml")
    authorization :=  "Bearer " + admin_token
    req.Header.Add("authorization", authorization)
    res, _ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println(string(body))

}


func Delete_app(appname string, pname string){
    var rcname string
    rcname = Get_rcname(appname, pname)
    podlist := Get_podlist(appname, pname, appname, rcname)
    var podname_str string

    //Delete deployconfig
    Delete_obj(appname, pname, appname, "deploymentconfigs")
 
    //Delete imagestream
    Delete_obj(appname, pname, appname, "imagestreams")

    //Delete service
    Delete_obj(appname, pname, appname, "services")
   
    // Delete pvc
    Delete_obj(appname, pname, appname, "persistentvolumeclaims")

    //Delete rc
    Delete_obj(appname, pname, rcname, "replicationcontrollers")

    //Delete pods
    for e := podlist.Front(); e != nil; e = e.Next() {
       podname := e.Value
       podname_str = fmt.Sprintf("%s", podname)
       Delete_obj(appname, pname, podname_str, "pods") 
    }
}


func Create_app(appname string, pname string, nodeport int, size int,replica int) {

   Create_obj(appname, pname, replica, nodeport, size, "persistentvolumeclaims")
   Create_obj(appname, pname, replica, nodeport, size, "imagestreams")
   Create_obj(appname, pname, replica, nodeport, size, "deploymentconfigs")
   Create_obj(appname, pname, replica ,nodeport, size, "services")
}




func main(){
    var nodeport int
    var size int
    var replica int

    appname := "test1"
    nodeport = 30008
    size = 10
    pname := "mysql"
    replica = 1
    
   Create_app(appname, pname, nodeport, size, replica)
 

    Delete_app("test1", "mysql")
     //Update_obj("test1", "mysql", 2, 30053, "test1", "port")
    //Update_obj("test1", "mysql", 2, 30054, "test1", "replica")
    //Get_rcname("owncloud", "owncloud")
    //Get_podlist("owncloud", "owncloud", "owncloud-1")   
    //Get_obj("owncloud", "owncloud", "owncloud", "services")
}

