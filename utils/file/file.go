package file

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func IsExist(path string) (bool ,error){
	_, err := os.Stat(path)
	if err == nil{
		return true , nil
	}
	if os.IsNotExist(err){
		return false ,err
	}
	return false , err
}

func ReadYaml(cpath string, cptr interface{}) error {
	bs, err := ioutil.ReadFile(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = yaml.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
}