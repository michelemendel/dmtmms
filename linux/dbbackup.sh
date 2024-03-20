#!/bin/bash

# chmod +x dbbackup.sh

# Set up the database file path
# db_path_prod="/home/michele/dmtmms/"
db_path="/Users/michelemendel/checkouts/dmtmms/devdb/"
db_file="dmtmms.db"
db_full_path="${db_path}${db_file}"

# Generate the backup file name with timestamp
backup_file="$(date +"%Y%m%dT%H%M")_dmtmms.db"
backup_full_path="${db_path}${backup_file}"

# Run the sqlite3 backup command
# This keeps a copy locally.
sqlite3 $db_full_path ".backup '$backup_full_path'"

# Copy the backup file to safe location
backup_copy_path="/Users/michelemendel/checkouts/"
cp $backup_full_path $backup_copy_path

echo "Backup complete: $backup_full_path, also copied to $backup_copy_path"