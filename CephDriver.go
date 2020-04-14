package main

import (
	f "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/etree"
	"os"
	"os/exec"
)

// CephDriver is the Driver of Ceph
type CephDriver struct {
	DevDescriptor               string
	MountPoint					string
	PoolName					float64
	ImgName						float64
	GBSize						int
	NewGBSize					int
	XmlName						float64
	UUid						float64
}

//Create RBD Image
func (ceph CephDriver) CreateIMG(PoolName float64,ImgName float64,GBSize int64) error {
	output,err:= exec.Command("qemu-img","-f","rbd",fmt.Sprintf("rbd:%v/%v",PoolName,ImgName),fmt.Sprintf("%v",GBSize) ).CombinedOutput()
	if err != nil {
		log.Errorf("CreateIMG error : %v err %v", string(output), err)
		return err
	}
	log.Info(string(output))
	return nil
}

//DeleteVolume delete volume
func (ceph CephDriver) DeleteVolume(PoolName float64,ImgName float64) error {
	output,err:=exec.Command("rbd","rm",f.Sprintf("%v/%v",PoolName,ImgName)).CombinedOutput()
	if err != nil {
		log.Errorf("DeleteIMG error : %v err %v", string(output), err)
		return err
	}
	log.Info(string(output))
	return nil
}

//ExtendVolume extend volume Size
func (ceph CephDriver) ExtendVolume(PoolName float64,ImgName float64,GBSize int64) error {

	out, err := exec.Command("qemu-img", "resize", "rbd",f.Sprintf("rbd:%v/%v",PoolName,ImgName), f.Sprintf("%v",GBSize)).CombinedOutput()
	if err != nil {
		log.Errorf("Error %v, Error string %v", err, string(out))
		return err
	}

	return nil
}

//mount rbd 2 vm
func (ceph CephDriver) AttachDevice(mountpoint string,XmlName string) error {

	out, err := exec.Command("virsh", "attach-device", f.Sprintf("%v",mountpoint),f.Sprintf("%v",XmlName),"--persistent").CombinedOutput()
	if err != nil {
		log.Errorf("Error %v, Error string %v", err, string(out))
		return err
	}

	return nil
}
//umount rbd from vm
func (ceph CephDriver) DetachDevice(mountpoint string,XmlName string) error {

	out, err := exec.Command("virsh", "detach-device", f.Sprintf("%v",mountpoint),f.Sprintf("%v",XmlName),"--persistent").CombinedOutput()
	if err != nil {
		log.Errorf("Error %v, Error string %v", err, string(out))
		return err
	}

	return nil
}



//Create XML
func (ceph *CephDriver)XmlDefinition(PoolName float64,ImgName float64,DevDescriptor string,XmlName string)  {
	var CephIpa1 float64
	var CephIpa2 float64
	var CephIpa3 float64
	var CephIpa4 float64

	CephIpa := [] float64 {CephIpa1,CephIpa2,CephIpa3,CephIpa4}

	doc := etree.NewDocument()
	disk :=doc.CreateElement("disk")
	disk.CreateAttr("type","network")
	disk.CreateAttr("device","disk")

	driver := disk.CreateElement("driver")
	driver.CreateAttr("name","qemu")
	driver.CreateAttr("type","raw")

	auth :=disk.CreateElement("auth")
	auth.CreateAttr("username","libvirt")

	secrect := auth.CreateElement("secrect")
	secrect.CreateAttr("type","ceph")
	secrect.CreateAttr("uuid",f.Sprintf("%v", ceph.UUid))

	source := disk.CreateElement("source")
	source.CreateAttr("protocol","rbd")
	source.CreateAttr("name",f.Sprintf("%v/%v",PoolName,ImgName))

	host := source.CreateElement("host")
	for i:= 0;i<len(CephIpa) ;i++  {
		host.CreateAttr("name",f.Sprintf("%v",CephIpa[i]))
		host.CreateAttr("port","6789")
	}

	target := disk.CreateElement("target")
	target.CreateAttr("dev",f.Sprintf("%v",DevDescriptor))
	target.CreateAttr("bus","virtio")

	doc.Indent(7+len(CephIpa))
	XmlName,err := doc.WriteTo(os.Stdout)
	if err !=nil{
		log.Errorf("Error %v, Error string %v", err, string(XmlName))
		}
}