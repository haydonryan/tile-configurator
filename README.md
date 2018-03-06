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

#### JSON Track (advanced customers):
Customer wants to manage entire tile configuration in JSON with pipelines and using commandline

##### Install
###### Customer has sandbox and can generate JSON from Ops Managaer
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
configure tile
save properties
(Note this doesn't work for X in ERT, keys, or secrets yet as the Ops Manger API won't return them)
###### Customer wants to start from JSON template / schema

##### Upgrade
Diff new tile to current tile to determine changes in json



#### YAML Track (easier):
 Customer wants to simplify to yaml
##### Install
###### Customer has sandbox and can generate JSON
###### Customer wants to start from YAML template

##### Upgrade
Diff new tile version to old to see what changes are needed (in yaml)
 With pipelines and using commandline



- task: unpack-pivotal-file
    config:
      platform: linux
      image_resource:
        type: docker-image
        source: {repository: busybox}
      inputs:
      - name: pivnet-gemfire-release
      run:
        path: /bin/sh
        args: ["-c", "find pivnet-gemfire-release/. -name '*.pivotal' | xargs -n1 unzip"]
      outputs:
      - name: releases