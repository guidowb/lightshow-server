kubectl create secret docker-registry gcr-registry-credentials \
  --docker-server=https://us.gcr.io \
  --docker-username=_json_key \
  --docker-email=user@example.com \
  --docker-password="$(cat ~/.secrets/k8s-gcr-auth-rw.json)"
