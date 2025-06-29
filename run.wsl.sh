#!/bin/bash
# WSL run script for the ROBO-9 platformer game
# This script sets the proper GOOS environment variable for running in WSL

echo "Starting ROBO-9 Platformer in WSL environment..."
GOOS=windows go run main.go
