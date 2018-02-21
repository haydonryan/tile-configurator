#!/bin/bash -e 
go test $(glide nv) -v
go test $(glide nv) -cover fmt