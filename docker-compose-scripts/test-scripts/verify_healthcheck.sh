#!/usr/bin/env bash
HEALTHCHECK_URL=$1

for i in {1..60}; do #timeout of 5 minutes

	status_code=$(curl --write-out '%{http_code}\n' --silent --output /dev/null $HEALTHCHECK_URL)
  echo "status code" $status_code
	if [ $status_code -eq 200 ]; then
		exit 0
	fi

	sleep 5
done

echo "Health check is not up"
exit 1
