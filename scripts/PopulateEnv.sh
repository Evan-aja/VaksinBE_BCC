#!/bin/bash

cd /var/app/current/;

if [[ -f ".env" ]];then
  echo "File env exist. Replacing with a new one";
  rm .env;
else
  echo "File env does not exist. Creating env file";
fi;

echo "DB_HOST=13.212.140.154:3306" >> .env;
echo "DB_NAME=intern_bcc_3" >> .env;
echo "DB_USER=admin" >> .env;
echo "DB_PASS=HnVXVx8rF4G3YjS3nKuQrKVS7apg4Vzt" >> .env;

echo "JWT_TOKEN=$(echo 'ThisIsASecretKey' | base64)" >> .env;

echo "TOKEN_G=It was one of the coldest winters, and many animals were dying because of the cold. The porcupines, realizing the situation, decided to group together to keep each other warm. This was a great way to protect themselves from the cold and keep each of them warm, but the quills of each one wounded their closest companions. After a while, they decided to distance themselves, but they too began to die due to cold. So they had to make a choice: either accept the quills of their companions or choose death. Wisely, they decided to go back to being together. They learned to live with a few wounds caused by their close relationship with their companions to receive the warmth of their togetherness. This way, they were able to survive." >> .env;
echo "CLIENT_ID=532005420615-cdqb43rk66g0r6h2n1kbjh0s2p5a2kuk.apps.googleusercontent.com" >> .env;
echo "CLIENT_SEC=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" >> .env;
