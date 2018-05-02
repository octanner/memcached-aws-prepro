## Synopsis

Docker image which runs cronjob for provisioning of memcached clusters on AWS ElastiCache

## Details
The cronjob makes sure that there are always 3 unclaimed small clusters, 1 unclaimed medium, and 1 unclaimed large


## Dependencies

1. "fmt"
2. "strings"
3. "database/sql"
4. "github.com/lib/pq"
5. "github.com/aws/aws-sdk-go/aws"
6. "github.com/aws/aws-sdk-go/aws/session"
7. "github.com/aws/aws-sdk-go/service/elasticache"
8. "github.com/nu7hatch/gouuid"
9. "os"



## Requirements
go

aws creds

## Runtime Environment Variables

1. LARGE_INSTANCE_TYPE
2. BROKER_DB
3. SMALL_INSTANCE_TYPE
4. MEDIUM_INSTANCE_TYPE
5. ELASTICACHE_SECURITY_GROUP
          


