package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"k8sToYml/myUtil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

var Datetime string
var RegisterResource = make(map[string]interface{})
var CmdPath string
var defaultYamlPath string

type Kind struct {
	Kind string `yaml:"kind"`
}
type ConfigMap struct {
	Kind       string            `yaml:"kind"`
	ApiVersion string            `yaml:"apiVersion"`
	Data       map[string]string `yaml:"data"`
	Metadata   Metadata          `yaml:"metadata"`
}
type Secret struct {
	Kind       string            `yaml:"kind"`
	ApiVersion string            `yaml:"apiVersion"`
	Data       map[string]string `yaml:"data"`
	Metadata   Metadata          `yaml:"metadata"`
	//StringData map[string]string `yaml:"stringData"`
	Type string `yaml:"type"`
}
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
	// Annotations map[string]string `yaml:"annotations"`
}
type CronJob struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type Job struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type ServiceAccount struct {
	Kind       string              `yaml:"kind"`
	ApiVersion string              `yaml:"apiVersion"`
	Secrets    []map[string]string `yaml:"secrets"`
	Metadata   Metadata            `yaml:"metadata"`
}
type PersistentVolumeClaim struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type PersistentVolume struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type Deployment struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type Service struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}
type StatefulSet struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Spec       interface{} `yaml:"spec"`
	Metadata   Metadata    `yaml:"metadata"`
}

func init() {
	RegisterResource["ConfigMap"] = &ConfigMap{}
	RegisterResource["Secret"] = &Secret{}
	RegisterResource["CronJob"] = &CronJob{}
	RegisterResource["ServiceAccount"] = &ServiceAccount{}
	RegisterResource["Job"] = &Job{}
	RegisterResource["PersistentVolumeClaim"] = &PersistentVolumeClaim{}
	RegisterResource["PersistentVolume"] = &PersistentVolume{}
	RegisterResource["Deployment"] = &Deployment{}
	RegisterResource["Service"] = &Service{}
	RegisterResource["StatefulSet"] = &StatefulSet{}
	// 复用
	RegisterResource["DaemonSet"] = &StatefulSet{}
	RegisterResource["ReplicaSet"] = &Deployment{}
	RegisterResource["HorizontalPodAutoscaler"] = &Deployment{}
	RegisterResource["Ingress"] = &Deployment{}
	CmdPath, _ = os.Getwd()
	Datetime = myUtil.DateTime()
	defaultYamlPath = CmdPath + "/<datetime>-<name>.yaml"
}

func main() {
	saveYamlPath := flag.String("s", defaultYamlPath, "save yaml path")
	oldYamlPath := flag.String("f", "", "k8s yaml file path")
	flag.Parse()
	if *oldYamlPath == "" {
		println("Please specify the k8a yaml path")
		return
	}
	if *saveYamlPath == defaultYamlPath {
		pathSplit := strings.Split(*oldYamlPath, "/")
		// defaultYamlPath = CmdPath + "/" + Datetime + "-" + pathSplit[len(pathSplit)-1]
		defaultYamlPath = CmdPath + "/" + "new-" + pathSplit[len(pathSplit)-1]
	} else {
		defaultYamlPath = *saveYamlPath
	}
	//var rs string
	var kind Kind
	strOldYaml, err := ioutil.ReadFile(*oldYamlPath)
	if err != nil {
		println(err)
	}
	// 把yaml形式的字符串解析成struct类型
	err = yaml.Unmarshal(strOldYaml, &kind)
	// fmt.Printf("kind: %s", kind.Kind)
	if err != nil {
		println("err:", err)
	}
	//rs = kind.Kind
	rsType1 := RegisterResource[kind.Kind]
	rr := reflect.TypeOf(rsType1).Elem()
	rsType := reflect.New(rr).Interface()
	//err = yaml.Unmarshal(config, &rsType) 反射 本身就是指针 去掉&
	err = yaml.Unmarshal(strOldYaml, rsType)
	if err != nil {
		println("err:", err)
	}
	//转换成yaml字符串类型
	builder := &strings.Builder{}
	encoder := yaml.NewEncoder(builder)
	// 设置缩进2个空格，注意要在Encode前面才生效
	encoder.SetIndent(2)
	err = encoder.Encode(rsType)
	myUtil.Check(err)
	err = encoder.Close()
	myUtil.Check(err)
	strYaml := builder.String()
	//strYaml, err := yaml.Marshal(rsType)
	err = ioutil.WriteFile(defaultYamlPath, []byte(strYaml), 0666)
	myUtil.Check(err)
	fmt.Printf("save success: %s\n", defaultYamlPath)
}
