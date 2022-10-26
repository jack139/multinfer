#!/bin/bash

LOG_PATH="/opt/ai/logs/"
LOG_PATH_B="/opt/ai/logs/backup/"

LOG_FILES=`ls $LOG_PATH*.log`

TO_DATE=`date +%Y%m%d`

rm -f $LOG_PATH_B/access* $LOG_PATH_B/backrun*

for file in $LOG_FILES
do
	cp $file $LOG_PATH_B`basename $file`.$TO_DATE
	cat /dev/null > $file
done

