# OAuth documentation

This documentation is focussed on the OAuth workflow, you should find all steps/information about the OAuth of Userstyles world(USw)
within this documentation and _most_ related coded is within the `handlers/oauthProvider/` folder.

_**Note:** currently headless OAuth isn't supported._

## Basic OAuth workflow

### 1 >> Let the user authorize with the given parameters
```
GET https://usestyles.world/api/oauth/authorize
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| state | string | Recommended | An unguessable random string. It is used to protect against [[https://en.wikipedia.org/wiki/Cross-site_request_forgery][cross-site request forgery]] attacks. |
| scope | string | No | A comma-delimited list of the [[./oauth.MD#scopes][scopes]], when not specified the default scopes will be used based on the application's setting. |


### 2 >> Users are sent to redirect_uri

_When the user accepts your request, USw redirects back to your site with a temporary code in a `code` parameter. This temporary code will expire after 10 minutes. When in the previous step the `state` parameter was provided, it will be sent back in the `state` parameter, When the states don't match, then a third party created the request, and you should abort the process._

```
POST https://userstyles.world/api/oauth/access_token
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| client_secret | string | Yes | The client secret that is generated when you've registered for an USw application. |
| code | string | Yes | The code that you've received as a parameter to Step 1. |
| state | string | No | The unguessable random string you've provided in Step 1. |


#### Response

Depending on the `Accept` header that was sent.
It will by default, send the `text/plain` format.

```
access_token=usw_IamAveryRandomAccesTokkeennAndImNotAfraid&token_type=Bearer
```

Otherwise when specified:
```json
Accept: application/json
{"access_token":"usw_IamAveryRandomAccesTokkeennAndImNotAfraid", "token_type":"Bearer"}
```

### 3 >> Using the access token

The access token will allow you to make requests to endpoints on behalf of the user.
The endpoints are specified [[./endpoints.MD][here]].

```
Authorization: TOKEN_TYPE ACCESS_TOKEN
GET https://userstyles.world/api/user/
```


## Specifc style workflow

This is an intresting workflow if you only care about to link with a specifc style.
This ensures the token is given only works with that specific style.

_Please note that that within the OAuth settings, the 'styles' scope still has to be set to actually use this._

### 1 >> Let the user authorize with the given parameters
```
GET https://usestyles.world/api/oauth/auth
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| state | string | Recommended | An unguessable random string. It is used to protect against [[https://en.wikipedia.org/wiki/Cross-site_request_forgery][cross-site request forgery]] attacks. |


### 2 >> Users are sent to redirect_uri

_When the user specified which style they want to link, USw redirects back to your site with a temporary code in a `code` parameter. This temporary code will expire after 10 minutes. When in the previous step the `state` parameter was provided, it will be sent back in the `state` parameter, When the states don't match, then a third party created the request, and you should abort the process. In this workflow it will also return a `style_id` parameter to specify which style has been chosen_

```
POST https://userstyles.world/api/oauth/token
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| client_secret | string | Yes | The client secret that is generated when you've registered for an USw application. |
| code | string | Yes | The code that you've received as a parameter to Step 1. |
| state | string | No | The unguessable random string you've provided in Step 1. |


#### Response

Depending on the `Accept` header that was sent.
It will by default, send the `text/plain` format.

```
access_token=usw_IamAveryRandomAccesTokkeennAndImNotAfraid&token_type=Bearer
```

Otherwise when specified:
```json
Accept: application/json
{"access_token":"usw_IamAveryRandomAccesTokkeennAndImNotAfraid", "token_type":"Bearer"}
```

### 3 >> Using the access token

The access token will allow you to make requests to style's endpoints on behalf of the user.
The endpoints are specified [[./endpoints.MD][here]].

```
Authorization: TOKEN_TYPE ACCESS_TOKEN
POST https://userstyles.world/api/style/1
```
data(unescaped):
```json
{
    "name": "New name!~"
}
```


## Scopes

_**Note:** Scopes may, later on, be split down into more specific scopes._

| Name | Description |
| :--- | :--- |
| `style` | Allow to add/edit/delete styles of the user. |
| `user` | Allow retrieving information of the user. |
