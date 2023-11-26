
import json
import logging

import requests
from flask import Flask, g
from flask_oidc import OpenIDConnect

logging.basicConfig(level=logging.DEBUG)

app = Flask(__name__)
app.config.update({
    'SECRET_KEY': 'SomethingNotEntirelySecret',
    'TESTING': True,
    'DEBUG': True,
    'OIDC_CLIENT_SECRETS': 'client_secrets.json',
    'OIDC_ID_TOKEN_COOKIE_SECURE': False,
    'OIDC_USER_INFO_ENABLED': True,
    'OIDC_OPENID_REALM': 'flask-demo',
    'OIDC_SCOPES': ['openid', 'email', 'profile'],
    'OIDC_INTROSPECTION_AUTH_METHOD': 'client_secret_post'
})

oidc = OpenIDConnect(app)


@app.route('/')
def public():
    if oidc.user_loggedin:
        return f"""Welcome, {oidc.user_getfield("preferred_username")}
                <ul>
                    <li><a href="/private">See private</a> </li>
                    <li><a href="/logout">Log out</a></li>
                </ul>
                """
    else:
        return 'Welcome, <a href="/private">Log in</a>'


@app.route('/private')
@oidc.require_login
def private():
    """Example for protected endpoint that extracts private information from the
    OpenID Connect id_token."""
    
    info = oidc.user_getinfo(['preferred_username', 'email', 'sub'])

    username = info.get('preferred_username')
    email = info.get('email')
    user_id = info.get('sub')

    return f"""Hi {username},<br><br>
            Your email is {email} and your user_id is {user_id}.\n\n
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="localhost:8080/realms/test/account?referrer=flask-app&referrer_uri=http://localhost:5000/private&">Account</a></li>
                <li><a href="/logout">Log out</a></li>
            </ul>"""


@app.route('/logout')
def logout():
    """Performs local logout by removing the session cookie."""
    oidc.logout()
    return 'You have been logged out! <a href="/">Return</a>'


if __name__ == '__main__':
    app.run()