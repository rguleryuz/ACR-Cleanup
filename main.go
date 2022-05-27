package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	AzureRegistryName = kingpin.Flag("azureregistryname", "ACR Name").Required().String() //req
	//SubscriptionName  = kingpin.Flag("subscriptionname", "Subscription Name").String()
	ImagestoKeep  = kingpin.Flag("imagestokeep", "Number of images to keep").Default("10").Int()
	EnableDelete  = kingpin.Flag("enabledelete", "Enable Delete").Default("no").String()
	Repository    = kingpin.Flag("repository", "Repository").String()
	deletedImages = 0
	repolist      []string
	imagetags     []string
	script        string
	ignoredtags   string
	tenantid      = os.Getenv("AZURE_TENANT_ID")
	spid          = os.Getenv("AZURE_CLIENT_ID")
	sppass        = os.Getenv("AZURE_CLIENT_SECRET")
)

func main() {
	checkOS()
	azlogin := exec.Command(script, "-c", "az login --service-principal -u "+spid+" -p "+sppass+" --tenant "+tenantid)
	loginout, _ := azlogin.Output()
	_ = loginout
	kingpin.Parse()
	//if SubscriptionName != nil {
	//	subscriptionName := exec.Command("bin/bash", "-c", "az account set --subscription"+*SubscriptionName)
	//
	//}
	if *Repository != "" {
		repolist = []string{"", "", *Repository}

	} else {
		getrepoout := exec.Command(script, "-c", "az acr repository list --name "+*AzureRegistryName+" --output table")
		getrepolist, _ := getrepoout.Output()
		repolist = strings.Fields(string(getrepolist))
	}
	for i := 2; i < len(repolist); i++ {
		index := repolist[i]
		fmt.Println("\nChecking repository: " + index)
		getimageout := exec.Command(script, "-c", "az acr repository show-tags --name "+*AzureRegistryName+" --repository "+index+" --output tsv --orderby time_desc"+ignoredtags)
		getimagetags, _ := getimageout.Output()
		imagetags = strings.Fields(string(getimagetags))
		fmt.Println("\n# Total images:", len(imagetags))
		fmt.Println("# Images to keep:", *ImagestoKeep)
		if len(imagetags) > *ImagestoKeep {
			fmt.Println("Deleting surplus images")
			for j := *ImagestoKeep; j < len(imagetags); j++ {
				imagename := *Repository + ":" + imagetags[j]
				deletedImages++
				if *EnableDelete == "yes" {
					fmt.Println("Deleting image: " + imagename)
					deleteImage := exec.Command(script, "-c", "az acr repository delete --name "+*AzureRegistryName+" --image "+imagename+" --yes")
					deleteImage.Run()
				} else {
					fmt.Println("dummy delete:" + imagename)
				}
			}
		} else {
			fmt.Println("No images to delete")
		}
	}
	fmt.Println("ACR cleanup Completed")
	fmt.Println("\n# Total images deleted:", deletedImages)
}

func checkOS() {
	if runtime.GOOS == "windows" {
		script = "pwsh"
		ignoredtags = " | Select-String -Pattern 'latest','dev','staging' -SimpleMatch -NotMatch"
	} else {
		script = "/bin/bash"
		ignoredtags = " | sed 's/\\(latest\\|dev\\|staging\\)//g'"
	}
}
