# Keycloak Benchmarking Guide

This guide facilitates the benchmarking of Keycloak through a Python Flask application integrated with Keycloak through OpenID Connect (OIDC). The setup is inspired by Thomas Darimont's example, accessible [here](https://gist.github.com/thomasdarimont/145dc9aa857b831ff2eff221b79d179a).

## Overview

- [Components](#components)
- [Set Up Keycloak](#set-up-keycloak)
- [Set Up Python Flask App](#set-up-python-flask-app)
- [Fill the Database](#fill-the-database)

## Components

- **flask-app**: Python Flask application
- **keycloak**: Keycloak installed in a Docker container
- **keycloak-fill-db-scripts**: Go scripts to populate the Keycloak database

## Set Up Keycloak

### 1. Launch Keycloak

```bash
docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:23.0.0 start-dev
```

- Note: Login credentials are set to `admin:admin`.
- Keycloak is now running at [http://localhost:8080](http://localhost:8080).

### 2. Create a Realm

![create-realm-image](images/create-realm.png)

### 3. Create a Client

#### General Settings

- Client ID: `flask-app`
  ![create-client-general-settings-image](images/create-client-general-settings.png)

#### Capability Config

- Client authentication: `On`
  ![create-client-capability-config-image](images/create-client-capability-config.png)

#### Login Settings

- Valid redirect URIs: `http://localhost:5000/*`
  ![create-client-login-settings-image](images/create-client-login-settings.png)

#### Copy Client Secret to Flask App

Copy the client secret from the client details page and paste it into the `client_secrets.json` file in the `flask-app` directory.
![copy-oidc-client-secret-image](images/copy-oidc-client-secret.png)

### 4. Adjust Token Lifetimes

By default the access token lifespan is set to 5 minutes and the SSO session idle timeout is set to 30 minutes. Make sure you change these settings in the `master` realm. To increase the token lifetimes, follow these steps:

- Set the SSO session idle timeout to 1 day:
  ![realm-settings-sessions-image](images/realm-settings-sessions.png)

- Set the access token lifespan to 1 day:
  ![realm-settings-tokens-image](images/realm-settings-tokens.png)

### 5. Create a User

- Username: `user`
- Email: `user@example.org`
- Email Verified: `On`
- First Name: `user`
- Last Name: `user`

![create-user-image](images/create-user.png)

#### Set User Password

- Make sure to turn off the `Temporary` switch before saving the password.

![create-user-credentials-image](images/create-user-credentials.png)

## Set Up Python Flask App

### 1. Create a Virtual Environment

```bash
mkvirtualenv keycloak-benchmarking
```

### 2. Install Dependencies

```bash
pip install flask flask_oidc
```

### 3. Run the Application

```bash
cd flask-app
python app.py
```

The Python Flask app is now running at [http://localhost:5000](http://localhost:5000).

## Fill the Database

### 1. Retrieve Admin User Access Token

```bash
curl -X POST \
  'http://localhost:8080/realms/master/protocol/openid-connect/token' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'username=admin' \
  --data-urlencode 'password=admin' \
  --data-urlencode 'grant_type=password' \
  --data-urlencode 'client_id=admin-cli'
```

### 2. Adjust Go Script Configurations

```go
const (
    baseUrl     = "http://localhost:8080"
    realm       = "your_realm"
    accessToken = "your_accessToken"
)
```

### 3. Install Go Modules

```bash
cd keycloak-fill-db-scripts
go get 
```

### 4. Execute Go Scripts

```bash
cd keycloak-fill-db-scripts 
go run ./cmd/create_groups <numberOfGroups>
```

Feel free to customize the configurations and adapt the scripts as needed for your testing environment.
