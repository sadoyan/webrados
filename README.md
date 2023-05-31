**WebRados**
---------

**WebRados** is a simple and high performance HTTP service for **CEPH** distributed file system.
It's for providing simple HTTP interface for **Ceph's** bare **RADOS** layer.
**WebRados** is not a replacement for **RadosGW** as is does not have all reach APIs and features of RadosGW(S3, Swift, etc ...), instead it stands for simplicity and performance.

WebRados relies on HTTP methods to interact with RADOS object, thus it can provide access to objects store in RADOS directly from internet browser .
The ide is to have web accessible storage for millions of relatively small files, which can be accessed from browser directly.

WebRados relies on C bindings of **Ceph**, so in order to run this program you need to install Ceph packages.
Running Ceph services on computer which hosts WebRados is not required, it's even better to have a dedicated server or pool of servers for running WebRados

### **Download and install**
---------

You can build WebRados from source or download precompiled binaries. If you already have installed Cephs packages and want to make things easy ,
just download te WebRados binary, make it executable, and you are ready to run .

Building from a source is also easy.

```shell
git clone  https://github.com/sadoyan/webrados.git
cd webrados
export GOROOT=/path/to/your/go
go mod tidy
go build .
./webrados -config config.yml 
```

### **Configuration**
---------

Example configuration file,  ```config.yml``` , with reasonable defaults is in root directory of the source tree.

```yaml
main:
  listen : 0.0.0.0:8080
  dispatchers : 20
  adminapikey : eyJleHAiOjAsImhhc2giOiI4ZWZjMGM0N2NlZGM5OWU1Y2EwM2QxNjIyMGU4ZGRkZDY2MjQ3NjI4ZTkyZTUwNThjMzc1MTY5Nzc2ODkzOTc4In0Q05m2RZV2JRQiW1Ho1yKm2iqqD4SS2qFibUAJcHKJ8Q
  dangerzone : yes
  readonly : no
  authread : yes
  authwrite : yes
  radoconns : 25
  logfile : no
  logpath: /opt/webrados.log
  allpools: no
  poollist:
    - bublics
    - donuts
    - images
  usersfile: users.txt
  authtype: signurl #apikey , basic, jwt, none
cache:
  shards: 1024
  lifewindow: 10
  cleanwindow: 1
  maxrntriesinwindow: 600000
  maxentrysize: 5000
  maxcachemb: 1024
monitoring:
  enabled : yes
  url:  127.0.0.1:9090
  user: admin
  pass: admin
```

### **API**
---------

| **Name**        | **Description**                                                  |
|-----------------|------------------------------------------------------------------|
| **Read File**   | HTTP ```GET:``` http://{BINDADDRESS}/{POOLNAME}/{FILENAME}       |
| **Upload File** | HTTP ```POST, PUT:``` http://{BINDADDRESS}/{POOLNAME}/{FILENAME} |
| **Remove File** | HTTP ```DELETE:``` http://{BINDADDRESS}/{POOLNAME}/{FILENAME}    |

Configuration file is pretty simple and intuitive.

### **Section main**
---------

| **Name**          | **Description**                                                                                                            |
|-------------------|----------------------------------------------------------------------------------------------------------------------------|
| **listen**        | IP port to bind.                                                                                                           |
| **dispatchers**   | Number of threads for webserver.                                                                                           |
| **serveruser**    | Static user.                                                                                                               |
| **serverpass**    | MD5 hash of password for static user.                                                                                      |                                                   | 
| **dangerzone**    | Enable HTTP DELETE method.                                                                                                 |
| **readonly**      | Enable readonly mode. If 'yes' only GET is allowed.                                                                        |
| **authread**      | Require authentication for GET only.                                                                                       |
| **authwrite**     | Require authentication for POST/PUT/DELETE.                                                                                |
| **radoconns**     | Number of connection to CEPH.                                                                                              |
| **logfile**       | Log to file, if 'no' logs are sent to stdout.                                                                              |
| **logpath**       | Path for log file.                                                                                                         |
| **allpools:**     | yes/no . If yes program will scan ceph and enable access via web to all pool.                                              | 
| **poollist:**     | Works only if **allpools** is set to **no**.                                                                               |  
| **usersfile**     | Path for file containing list of users with `username passwordhash` format separated by new line.                          |
| **authtype**      | Authentication methods. ***apikey*** (X-API-KEY Header) , **basic** (HTTP Basic Auth), **jwt** (https://jwt.io/), **none** |

### **Section cache**
---------

| **Name**               | **Description**                                                                                   |
|------------------------|---------------------------------------------------------------------------------------------------|
| **shards**             | Number of shards (must be a power of 2)                                                           |
| **lifewindow**         | Time after which entry can be evicted                                                             |
| **cleanwindow**        | Interval between removing expired entries (clean up). If set to <= 0 then no action is performed. |
| **maxrntriesinwindow** | rps * lifeWindow, used only in initial memory allocation                                          |
| **maxentrysize**       | max entry size in bytes, used only in initial memory allocation                                   |
| **maxcachemb**         | Cache will not allocate more memory than this limit, value in MB.  0 value means no size limit    |

### **Section monitoring**
---------

| **Name**    | **Description**                               |
|-------------|-----------------------------------------------|
| **enabled** | Enable/Disable monitoring.                    |
| **url**     | IP address and port for minitoring interface. |
| **user**    | Monitoring user.                              |
| **pass**    | Password for monitoring user.                 |


### **Authentication and users.txt file**

Webrados can dynamically update users and API keys from ```users.txt``` file .
If you are using Basic Auth, ```users.txt``` should contain user and md5hash of password divided by space in each line.  

On Linux systems ```echo -n SecretPaSs | md5sum |awk '{print $1}'```  will output md5hash for using it as password in ```users.txt``` file

If you are using API keys ```users.txt``` should contain these keys separated by new line. 

Webrados will periodically read ```uesrs.txt``` file and automatically update users in memory.

If you are using JWT Authentication, you should set the value of your JWT Secret as **JWTSECRET** OS enviroment. 

```
export JWTSECRET='Super$ecter123765@'
```

You can pass JWT token as header `Authorization` . 

```
curl --header "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/bublics/nice.jpeg
```
Or as URL parameter. 
```
curl http://127.0.0.1:8080/bublics/nice.jpeg?token=$TOKEN
```

### **Large files**

In order to be able to store large file in RADOS directly files needs to be split to smaller chunks.
WebRados will automatically set maximum chunk size to  **OSDMaxObjectSize** of Ceph and split files in accordance to that.

### **Special commands**

**HTTP GET** http://{BINDADDRESS}/{POOLNAME}/{FILENAME}?info 
Returns information about requested file in json format.

```
curl -s  http://ceph1:8080/bublics/katana.mp4?info
```

```json
{
  "name": "katana.mp4",
  "pool": "bublics",
  "size": 471861144,
  "parts": 11
}
```

**HTTP DELETE** http://{BINDADDRESS}/{POOLNAME}/{FILENAME}?cache
Removes entry of given file from metadata cache

```curl -XDELETE  http://ceph1:8080/bublics/katana.mp4?cache```

**Admin Endpoint**

http://{BINDADDRESS}/.admin is special endpoint for performing administration tasks. 
This endpoint is always secured with special API-KEY solely for admin tasks. `config.yml[main][adminapikey]`. 
Key in existing config is just an example, please change it.

Admin commands are working on http POST, PUT, GET methods. Here are examples. 

```shell
curl -H "X-API-KEY: $APIKEY" 'http://ceph1:8080/.admin?purgecachestats'
```
Purge Statistics for local cache
```shell
curl -H "X-API-KEY: $APIKEY" 'http://ceph1:8080/.admin?purgecache'
```
Empty local cache
```shell
curl -H "X-API-KEY: $APIKEY" -d '{"url": "http://ceph1:8080","exp": 1685365532}'  'http://ceph1:8080/.admin?genjwt'
```
Generate JWT token with TTL till `exp` 

```shell
curl -s -XPOST -H "X-API-KEY: $APIKEY" --data-binary @/tmp/data.json 'http://ceph1:8080/.admin?sign'
```
Sends URLs for signing. Returns json with url and signed url pairs. 
URL signing will use `JWTSECRET` environment variable as a secret.  
data.json should be a json file with url(**string**) as key with value of TTL(**int**). TTL is seconds after **now**. 
Example request json: 
```json
{
  "http://ceph1:8080/bublics/1.jpeg": 180,
  "http://ceph1:8080/bublics/2.jpeg": 120,
  "http://ceph1:8080/bublics/3.jpeg": 90,
  "http://ceph1:8080/bublics/4.jpeg": 150
}
```
Example response :
```json
{
  "http://ceph1:8080/bublics/1.jpeg": "http://ceph1:8080/bublics/1.jpeg?expiry=1685441910&signature=G0L9EgRRYavGX8TRqLsH-ifIRm3__mFkR95wR3r4TyQ",
  "http://ceph1:8080/bublics/2.jpeg": "http://ceph1:8080/bublics/2.jpeg?expiry=1685441850&signature=psSGxImNZdzuSixB8dbaUhlFO81Mj02q6zLUrB_Rajs",
  "http://ceph1:8080/bublics/3.jpeg": "http://ceph1:8080/bublics/3.jpeg?expiry=1685441820&signature=Fr9OI0d_UMnAGZIPdt_tTl6QywH5Z0R_dnyfEgyxH3M",
  "http://ceph1:8080/bublics/4.jpeg": "http://ceph1:8080/bublics/4.jpeg?expiry=1685441880&signature=3vG0YGxSRlxXo7xWyi5l7Otw8Y5WZteq7mfUUq26lKE"
}

```
