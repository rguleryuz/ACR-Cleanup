# Introduction 
Azure Container Registry Image cleanup script written in GO


# How to Use
## ENV Variables

AZURE_TENANT_ID - Azure Tenant id of service principal

AZURE_CLIENT_ID - Azure Client id of service principal 

AZURE_CLIENT_SECRET - Azure Client password of service principal


## ARGS

azureregistryname (Required) - Name of the ACR Resource

repository(optional) - specify the repository to scan for surplus images to delete. If omitted all repositories will be be scanned within the specified container registry eg. smartpulse/management

imagestokeep (optional) - How many recent images you want to keep (Default = 10, if not stated)

enabledelete (optional) - Enable actual deletion of images instead of just scanning for surplus images (default = no, change to "yes" to delete images)

# Example

1. go run main.go -azureregistryname yourregistryname

In this case, script will default to scanning for images to delete (ImagestoKeep = 10 and EnableDelete = "no")

2. go run main.go -azureregistryname yourregistryname -enabledelete yes

In this case, script will delete any surplus images above the default of 10 images

3. go run main.go -azureregistryname yourregistryname -enabledelete yes -imagestokeep 20

In this case, script will delete any surplus images above 20 images per repository

4. go run main.go -azureregistryname yourregistryname -enabledelete yes -repository yourrepositoryname

In this case, script will delete any surplus images within the "yourrepositoryname" repository
