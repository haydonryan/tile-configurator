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




