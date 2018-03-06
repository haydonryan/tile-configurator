# tile-configurator 

### Purpose
This is a hacked together tool that will read in a yaml file and orchestrate the OM tool (https://github.com/pivotal-cf/om) (and therefore ops manager) to install and update tiles.


### Issues / Feature requests
If you have an issue with this tool - please submit a github issue.  I'm using Pivotal Tracker to manage the backlog for this project, so please contact me directly for access (https://www.pivotaltracker.com/n/projects/2151323).

### Benefits
- Configuration is in YAML rather than jSON - easier to configure, read and maintain.
- Externalization of configuration (cloud native operations)
- The tool will apply properties one at a time for individual properties, or in groups/collections for properties that need to be applied together.  This improves the feedback that the user recieves.
- Becasue it's a go app, we can add retreiving of credentials without operators seeing them (or being able to check the environment) (not yet)

### To Do
- Injest:
  - Read tile metadata into tool to provide help.
  - Externalise property name lookup dictionary.
- Configure:
  - Add flag to apply everything in one hit rather than individually (faster but less feedback for errors)
  - Add output to json 
  - Improve the application by adding tests, and updating the cli according to https://blog.alexellis.io/5-keys-to-a-killer-go-cli/


### Workflows:
Tile-configurator has been built to support the two main workflows that we see customers want to do.  

#### Starting with a manual install of PCF.
By doing a manual install at the start it allows inexperienced cloud operators understand the process of configuring Pivotal Cloud Foundry and it's components. 

#### Starting with an automated deployment of PCF.
Starting with an automated deploy of PCF allows the solutions team to biuld configuration applicable to the customer before going on site.  It also allows temporary labs to be created and full Infrastructure as Code.

However the end goal is the same - automated PCF deploys.


### Customer has sandbox deployed and can generate JSON from Ops Managaer

Stage tile
```
$ om-linux -t $TARGET -u $USERNAME -p $PASSWORD -k staged-products
+---------+-----------------+
|  NAME   |     VERSION     |
+---------+-----------------+
| p-bosh  | 1.10.12.0       |
| cf      | 1.10.20-build.6 |
| p-mysql | 1.9.18          |
+---------+-----------------+
```

Save properties
```

# List all staged products 
$ om-linux -t $TARGET -u $USERNAME -p $PASSWORD -k curl -p /api/v0/staged/products --silent
[
  {
    "installation_name": "p-bosh",
    "guid": "p-bosh-e264d97fa75e1646f473",
    "type": "p-bosh",
    "product_version": "1.10.12.0"
  },
  {
    "installation_name": "cf-0fc76391fde7e5f1ad58",
    "guid": "cf-0fc76391fde7e5f1ad58",
    "type": "cf",
    "product_version": "1.10.20-build.6"
  },
  {
    "installation_name": "p-mysql-036173c998d92c2ce8ad",
    "guid": "p-mysql-036173c998d92c2ce8ad",
    "type": "p-mysql",
    "product_version": "1.9.18"
  }
]

# Use the product GUID to get the properties and redirect this into a file
$ om-linux -t $TARGET -u $USERNAME -p $PASSWORD -k curl -p /api/v0/staged/products/p-mysql-036173c998d92c2ce8ad/properties --silent > mysql-properties.json

```

Now that we have the properties in a json file we can ingest it to a yaml with comments.  

  

```
# The -a annotates the yaml with comments generated (currently) from help.yml. 
# The -s simplifies the keys used by doing a dictionary lookup.

$ ./tile-configurator ingest -i ~/workspace/dojostarter/tiles/p-mysql-1.10.12.json -s -a > output.yml

```

The contents of the output.yml will look something like this:
```
# ===================================
# MySQL v1.9 GENERATED Parameters
# ===================================

# ===================================
# MySQL Server parameters
# ===================================

.cf-mysql-broker.bind_hostname:                                       #- hostname for mysql broker
.cf-mysql-broker.quota_enforcer_pause: 30                             #- Configure how many seconds the Quota Enforcer pauses between polls. Advanced configuration, please read the documentation before modifying. (default 30)
.mysql.allow_local_infile: true
.mysql.allow_remote_admin_access: false
.mysql.binlog_expire_days: 7
.mysql.cli_history: true
.mysql.cluster_name: cf-mariadb-galera-cluster                        #- Cluster name of the MySQL cluster - do not change when upgrading!
.mysql.cluster_probe_timeout: 10
.mysql.innodb_large_prefix_enabled: true
.mysql.innodb_strict_mode: true
.mysql.max_connections: 1500
.mysql.metrics_polling_frequency: 30                                  #- Select the polling interval for MySQL metrics in seconds
.mysql.mysql_start_timeout: 60
.mysql.roadmin_password: ***
.mysql.skip_name_resolve: true
.mysql.table_definition_cache: 8192
.mysql.table_open_cache: 2000
.mysql.tmp_table_size: 33554432
# ===================================
# Enable Backups
# ===================================
# (Note: If you choose disable then you also need to set backup prepare node instances to 0 in resources)

backup: disable                                                       #- Choose (disable  | enable)
backup_masters: true                                                  #- Each node is a duplicate master. This option makes unique backups from each master, rather than from a single instance.
backup_cron_schedule:                                                 #- Cron Schedule (See http://godoc.org/github.com/robfig/cron)
# ==================================
# Set Backup Destination
# ===================================
# (Note enable means enable s3)
export_backups: disable                                               #- Choose (disable  | enable | azure | gcs | scp)
azure_backup_base_url:                                                #- URL of Azure BlobStore
azure_backup_container:                                               #- Azure Container Name
azure_backup_stir:                                                    #- Azure Container Name
azure_backup_storage_account_key:                                     #- Azure Storage Access Key
azure_backup_storage_account_name:                                    #- Azure Storage Account Name

s3_backup_access_key:                                                 #- S3 Access Key
s3_backup_bucket_name:                                                #- Bucket Name
s3_backup_bucket_path:                                                #- Bucket Path
s3_backup_endpoint:                                                   #- S3 Endpoint 
s3_backup_region:                                                     #- S3 region (If using AWS S3, this field is required for any non us-east-1 regions)
s3_backup_access_key:                                                 #- S3 secret key
```


### Comparing configurations
There are many scenarios that you want to work out what has changed between configurations. This may include things like:
- Diff of what has been configured from a staged tile
- Diff between two environments ie sandbox and development
- Diff between what's in source control and has been configured in Ops Manager
- Diff between what version N of the tile contains, and N+1

Tile Configurator includes a diff feature.  It will report the keys that have been added to a tile. Note that to find out the properties that have been removed, you just swap the parameters.  This was done so that the outputs could be captured and annotated with comments.

```
$ ./tile-configurator diff -b fixtures/p-mysql-1.9.10-unconfigured.json -c fixtures/p-mysql-1.9.18-unconfigured.json  |jq
{
  "properties": {
    ".cf-mysql-broker.allow_table_locks": {
      "configurable": true,
      "credential": false,
      "optional": false,
      "type": "boolean",
      "value": true
    },
    ".mysql.mysql_backup_server_certificate": {
      "configurable": false,
      "credential": true,
      "optional": false,
      "type": "rsa_cert_credentials",
      "value": {
        "private_key_pem": "***"
      }
    },
    ".properties.innodb_flush_log_at_trx_commit": {
      "configurable": true,
      "credential": false,
      "optional": false,
      "type": "selector",
      "value": "two"
    },
    ".properties.syslog.enabled.protocol": {
      "configurable": true,
      "credential": false,
      "optional": true,
      "options": [
        {
          "label": "TCP protocol",
          "value": "tcp"
        },
        {
          "label": "UDP protocol",
          "value": "udp"
        },
        {
          "label": "RELP protocol",
          "value": "relp"
        }
      ],
      "type": "dropdown_select",
      "value": "tcp"
    }
  }
}

```

This file can then be injested back into yaml to have comments and easy to read keys as above.


### Applying configuration to Ops Manager

The final function of this tool is to provide good feedback when applying configurationt to Ops Manager.  

The goals for this component:
 - Configuration applied one parameter at a time or in groups where required.
 - Hide details when things go right - ie applying properties.
 - Provide clear communication to the worker when things to wrong
 - Support running in a pipeline / command line for fast feedback and therefore iterations.

Example showing error when applying a group of properties:
 ```
 $ ./tile-configurator config -i fixtures/properties.yml -t p-mysql -u <<REDACTED>> -p <<REDACTED>> --url opsmgr.sandbox.lab.hnrglobal.com
Converting syslog_forwarder to .properties.syslog
Applying setting .properties.plan_collection to tile p-mysql......Done.
Applying setting .properties.syslog to tile p-mysql......Done.
Applying Group optional-protections to tile p-mysql.....Call to OM returned error:
exit status 1
configuring product...
setting properties
could not execute "configure-product": failed to configure product: request failed: unexpected response:
HTTP/1.1 422 Unprocessable Entity
Transfer-Encoding: chunked
Cache-Control: no-cache, no-store
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Date: Tue, 06 Mar 2018 15:34:30 GMT
Expires: Fri, 01 Jan 1990 00:00:00 GMT
Pragma: no-cache
Server: nginx/1.4.6 (Ubuntu)
Strict-Transport-Security: max-age=31536000
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Request-Id: 82d9f180-a017-4b67-806c-62177637577a
X-Runtime: 1.041311
X-Xss-Protection: 1; mode=block

49
{"errors":{"optional_protections_recipient_email":["is not a property"]}}
0

```


