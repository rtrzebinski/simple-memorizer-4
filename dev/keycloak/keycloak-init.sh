#!/bin/bash
set -e

KEYCLOAK_URL="http://localhost:8180"
ADMIN_USER="admin"
ADMIN_PASS="change_me"
REALM_NAME="realm-dev"
CLIENT_ID="client-id-dev"
CLIENT_SECRET="client-secret-dev"

until curl -s "$KEYCLOAK_URL/admin/realms" > /dev/null; do
  echo "[Keycloak] Waiting for admin endpoint..."
  sleep 2
done
echo "[Keycloak] Keycloak is responding!"

echo "[Keycloak] Patching Keycloak master realm SSL..."
docker-compose -f dev/docker-compose.yml exec sm-keycloak bash -lc "\
  /opt/keycloak/bin/kcadm.sh config credentials --server $KEYCLOAK_URL --realm master --user $ADMIN_USER --password $ADMIN_PASS && \
  /opt/keycloak/bin/kcadm.sh update realms/master -s sslRequired=NONE
" > /dev/null 2>&1

echo '[Keycloak] Fetching admin token...'
ADMIN_TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASS" \
  -d "grant_type=password" \
  -d "client_id=admin-cli" \
  | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')

echo "[Keycloak] Creating realm '$REALM_NAME'..."
curl -s -X POST "$KEYCLOAK_URL/admin/realms" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"realm\":\"$REALM_NAME\",
    \"enabled\":true,
    \"sslRequired\":\"NONE\"
  }" || true

docker-compose -f dev/docker-compose.yml exec sm-keycloak bash -lc "\
  /opt/keycloak/bin/kcadm.sh config credentials --server $KEYCLOAK_URL --realm master --user $ADMIN_USER --password $ADMIN_PASS && \
  /opt/keycloak/bin/kcadm.sh update realms/$REALM_NAME \
    -s registrationEmailAsUsername=true \
    -s duplicateEmailsAllowed=false \
    -s editUsernameAllowed=false
" > /dev/null 2>&1

echo "[Keycloak] Enabling refresh rotation, reuse detection, and token lifetimes on realm..."
docker-compose -f dev/docker-compose.yml exec sm-keycloak bash -lc "\
  /opt/keycloak/bin/kcadm.sh config credentials \
    --server $KEYCLOAK_URL \
    --realm master \
    --user $ADMIN_USER \
    --password $ADMIN_PASS && \
  /opt/keycloak/bin/kcadm.sh update realms/$REALM_NAME \
    -s revokeRefreshToken=true \
    -s refreshTokenMaxReuse=0 \
    -s accessTokenLifespan=600 \
    -s ssoSessionIdleTimeout=2592000 \
    -s ssoSessionMaxLifespan=2592000
" > /dev/null 2>&1

echo "[Keycloak] Disabling default required actions..."
docker-compose -f dev/docker-compose.yml exec sm-keycloak bash -lc "\
  /opt/keycloak/bin/kcadm.sh config credentials --server $KEYCLOAK_URL --realm master --user $ADMIN_USER --password $ADMIN_PASS && \
  /opt/keycloak/bin/kcadm.sh update realms/$REALM_NAME -s verifyEmail=false && \
  for a in VERIFY_EMAIL UPDATE_PASSWORD UPDATE_PROFILE CONFIGURE_TOTP TERMS_AND_CONDITIONS; do \
    /opt/keycloak/bin/kcadm.sh update authentication/required-actions/\$a -r $REALM_NAME -s defaultAction=false -s enabled=false || true; \
  done
" > /dev/null 2>&1

echo "[Keycloak] Creating confidential client '$CLIENT_ID' with predefined secret..."
curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"clientId\":\"$CLIENT_ID\",
    \"secret\":\"$CLIENT_SECRET\",
    \"publicClient\":false,
    \"directAccessGrantsEnabled\":true,
    \"standardFlowEnabled\":true,
    \"serviceAccountsEnabled\":true,
    \"consentRequired\":false
  }" || true
echo "[Keycloak] Client secret for $CLIENT_ID: $CLIENT_SECRET"

echo "[Keycloak] Waiting for service account to appear..."
CLIENT_UUID=""
SERVICE_ACCOUNT_ID=""
for i in {1..30}; do
  CLIENT_UUID=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=$CLIENT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | sed -n 's/.*\"id\":\"\([^\"]*\)\".*/\1/p')
  if [ -n "$CLIENT_UUID" ]; then
    SERVICE_ACCOUNT_ID=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$CLIENT_UUID/service-account-user" \
      -H "Authorization: Bearer $ADMIN_TOKEN" | sed -n 's/.*\"id\":\"\([^\"]*\)\".*/\1/p')
    [ -n "$SERVICE_ACCOUNT_ID" ] && break
  fi
  sleep 1
done
if [ -z "$SERVICE_ACCOUNT_ID" ]; then
  echo "[Keycloak] ERROR: Service account did not appear within 30 seconds"
  exit 1
fi
echo "[Keycloak] Service account ready: $SERVICE_ACCOUNT_ID"

echo "[Keycloak] Granting realm-management client roles to service account..."
REALM_MGMT_UUID=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients?clientId=realm-management" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | sed -n 's/.*\"id\":\"\([^\"]*\)\".*/\1/p')

MANAGE_USERS_ROLE_JSON=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$REALM_MGMT_UUID/roles/manage-users" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
MANAGE_USERS_ROLE_ID=$(echo "[Keycloak] $MANAGE_USERS_ROLE_JSON" | sed -n 's/.*\"id\":\"\([^\"]*\)\".*/\1/p')

VIEW_USERS_ROLE_JSON=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients/$REALM_MGMT_UUID/roles/view-users" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
VIEW_USERS_ROLE_ID=$(echo "[Keycloak] $VIEW_USERS_ROLE_JSON" | sed -n 's/.*\"id\":\"\([^\"]*\)\".*/\1/p')

ASSIGN_PAYLOAD="["
if [ -n "$MANAGE_USERS_ROLE_ID" ]; then
  ASSIGN_PAYLOAD="$ASSIGN_PAYLOAD{\"id\":\"$MANAGE_USERS_ROLE_ID\",\"name\":\"manage-users\",\"clientRole\":true,\"containerId\":\"$REALM_MGMT_UUID\"}"
fi
if [ -n "$VIEW_USERS_ROLE_ID" ]; then
  [ "$ASSIGN_PAYLOAD" != "[" ] && ASSIGN_PAYLOAD="$ASSIGN_PAYLOAD,"
  ASSIGN_PAYLOAD="$ASSIGN_PAYLOAD{\"id\":\"$VIEW_USERS_ROLE_ID\",\"name\":\"view-users\",\"clientRole\":true,\"containerId\":\"$REALM_MGMT_UUID\"}"
fi
ASSIGN_PAYLOAD="$ASSIGN_PAYLOAD]"

if [ "$ASSIGN_PAYLOAD" != "[]" ]; then
  curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$SERVICE_ACCOUNT_ID/role-mappings/clients/$REALM_MGMT_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$ASSIGN_PAYLOAD" > /dev/null || true
  echo "[Keycloak] Assigned realm-management client roles to service account."
else
  echo "[Keycloak] WARNING: Could not resolve realm-management roles; skipping assignment."
fi

echo "[Keycloak] Realm '$REALM_NAME', client '$CLIENT_ID' created!"
echo "[Keycloak] Keycloak UI -> $KEYCLOAK_URL  (admin:change_me)"
