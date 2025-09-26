#!/bin/bash
gocyclo . | awk '{print $1}'

