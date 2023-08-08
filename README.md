# Configure your IP to the private ACR to perform the Docker operation

## When you create a private ACR that time you are not allow to perform any operatio in that ACR so you need to register you IP with the ACR or you need to run the code in the same VNet where you endpoint is created for private ACR. So in this case we are using the ip configuration option to configure the private ACR.


### Follow the below steps to run the terratest code:


You need to define these values in your code before running:-

            	resourceGroup := "<BLOG>"  // where the ACR created
                acrName := "<testacrmk2>"  // ACR name
                acrURL := "<testacrmk2.azurecr.io>" // ACR login server URL
                acrPassword := "< >" // ACR password
                imageNameLocal := "hello-world:latest" // image name

Step 1:- Run the go initialization command:

            go mod init < name >

Step 2:- Run the tidy command to install the packages:-

            go mod tidy

Step 3:- Run the test command:-

            go test -v

