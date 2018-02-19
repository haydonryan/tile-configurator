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
  - Make an injestor to create starting templates from the API
  - Add diff to fresh staged properties to allow creation of a template that is only "what's changed"
  - Externalise property name lookup help.
- Configure:
  - Add flag to apply everything in one hit rather than individually (faster but less feedback for errors)
  - Add output to json 
  - Improve the application by adding tests, and updating the cli according to https://blog.alexellis.io/5-keys-to-a-killer-go-cli/


### Workflows:

#### JSON Track (advanced customers):
Customer wants to manage entire tile configuration in JSON with pipelines and using commandline

##### Install
###### Customer has sandbox and can generate JSON from Ops Managaer
stage tile
save properties
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
