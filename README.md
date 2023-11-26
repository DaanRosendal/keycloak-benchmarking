# Keycloak Benchmarking

This guide provides a simple test environment for a Python Flask application integrated with Keycloak through OpenID Connect (OIDC). The setup is inspired by the somewhat outdated example by Thomas Darimont [here](https://gist.github.com/thomasdarimont/145dc9aa857b831ff2eff221b79d179a).

## Components

- **flask-app**: Python Flask application
- **keycloak-23.0.0**: Keycloak bare metal installation
- **keycloak-fill-db-scripts**: Go scripts to populate the Keycloak database

## Set Up Python Flask App

1. Create a new virtual environment:

```bash
mkvirtualenv keycloak-benchmarking
```

2. Install Flask and Flask-OIDC:

```bash
pip install flask flask_oidc
```

3. Run the application:

```bash
python flask-app/app.py
```

The Python Flask app is now running at [http://localhost:5000](http://localhost:5000).

## Set Up Keycloak

1. Launch Keycloak:

```bash
keycloak-23.0.0/bin/kc.sh start-dev
```

Keycloak is now running at [http://localhost:8080](http://localhost:8080).

2. Follow the steps in Keycloak's official ["Getting Started" documentation](https://www.keycloak.org/getting-started/getting-started-zip) for additional setup.

3. When creating the client, ensure that `Valid Redirect URIs` is set to `http://localhost:5000/*`.

![image](https://github.com/DaanRosendal/keycloak-benchmarking/assets/32291500/f8613e2e-f4e5-45b9-a53a-1246930c05dd)

## Fill the Database

1. Adjust the configurations at the top of the file, e.g., for creating groups:

```go
const (
    baseUrl     = "http://localhost:8080"
    realm       = "your_realm"
    accessToken = "your_accessToken"
)
```

2. Execute the scripts:

```bash
cd keycloak-fill-db-scripts
go run ./cmd/create_groups <numberOfGroups>
go run ./cmd/add_user_to_groups
```

Feel free to customize the configurations and adapt the scripts as needed for your testing environment.
