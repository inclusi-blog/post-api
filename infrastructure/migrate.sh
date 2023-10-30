#!/usr/bin/env sh
SECRET_JSON=$(aws secretsmanager get-secret-value --secret-id $SECRET_KEY --query 'SecretString' --output text)
ENV_VARS=$(echo $SECRET_JSON | jq -r 'to_entries[] | "\(.key)=\(.value)"')

# Set the environment variables
eval $ENV_VARS
export STORY_PASSWORD=$STORY_PASSWORD

if [ $? -eq 0 ]; then
    /home/gola/flyway -url=jdbc:postgresql://"${DB_HOST}":"${DB_PORT}"/"${DB_NAME}" -schemas="${DB_NAME}" -user="${DB_USER}" -password="${STORY_PASSWORD}" -connectRetries=60 -mixed=true migrate
else
    echo "Unable to substitute credentials"
    exit 1
fi
