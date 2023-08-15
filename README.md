# Golang-Ldap-API
LDAP actions over REST-API. 

# Functions this API can Perform:
Following are the Actions this API can perform.
- Create User in AD
- Create Security Group in AD
- Remove Security Group from AD
- Add User to Security Group in AD
- Remove User from Security Group in AD
- Add Bulk Users to a Security Group in AD
- Remove Bulk User from a Security Group in AD
- Check if User Exist in a Security Group.

# Docs:
http://localhost/docs

# Config:
- Bind_Address (Port you want this API to Live on Example :8080)
- LDAP_SERVER (Your Main AD LDAP Address to connect Example: ldaps://my-main-ad.mydomain.com)
- LDAP_SERVER2 (This is for Backup AD, If you have multiple ADs then enter another AD URL or if you don't have any other AD then just enter the first one here too.)
- LDAP_USER (LDAP Admin User to bind with AD, Note: In ldapfunc/conn.go OU for this user is Hardcoded which is "Admin", if your Ldap_user OU is different then change that OU from the file.)
- LDAP_PASSWORD (Password if the LDAP_USER)
- Domain_First (The First part of your AD Domain Example: AD-Domain is my-main-ad.mydomain.com, mydomain is Domain_First)
- Domain_Last (The Last part of your AD Domain Example: AD-Domain is my-main-ad.mydomain.com, com is Domain_Last)
- API_TOKEN (Token use in the Authentication calls from client. You can give multiple tokens in comma seperate string. Example: token1,token2,token3)

This all params get from environment variable so these should be define in the env. Or you can directly input this in Util/config.go

# Docker File and Docker-Compose:

Docker File and Docker Compose is also attached in the repository. Just set params in Docker file and compose and docker-compose up and you are ready to go.
You also need to provide your domain cert and key in order to use this API in HTTPS mode. You can define crt and key path in the nginx config.

# Kubernetes Deployment:
Kubernetes Deployment is also present in the repository. You just need to create an image of this program using docker file then place that image in Docker HUB and create secrets of your docker hub in the kubernetes and update that secret name in the Deployment file.
