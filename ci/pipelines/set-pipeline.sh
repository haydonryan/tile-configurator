#!/bin/bash

# Must have already logged in


fly -t home set-pipeline -p ops-stage-configure -c load-tile-stage-configure.yml --load-vars-from=params.yml
