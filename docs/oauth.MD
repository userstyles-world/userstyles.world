# OAuth documentation

This documentation is focussed on the OAuth workflow, you should find all steps/information about the OAuth of Userstyles world(USw)
within this documentation and _all_ related coded is within the `oauth_provider/` folder.

_**Note:** currently headless OAuth isn't supported._

## Basic OAuth workflow

### 1 >> Let the user authorize with the given parameters  
```
GET https://usestyles.world/oauth/authorize
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| state | string | Recommended | An unguessable random string. It is used to protect against [[https://en.wikipedia.org/wiki/Cross-site_request_forgery][cross-site request forgery]] attacks. |
| redirect_uri | string | No | Specify which URL should be used after the authorization is finished. When not specified, default value within the application would be used. |
| scope | string | No | A comma-delimited list of the [[./oauth.MD#scopes][scopes]], when not specified the default scopes will be used based on the application's setting. |


### 2 >> Users are sent to redirect_uri  

_When the user accepts your request, USw redirects back to your site with a temporary code in a `code` parameter. This temporary code will expire after 10 minutes. When in the previous step the `state` parameter was provided, it will be sent back in the `state` parameter, When the states don't match, then a third party created the request, and you should abort the process._

```
POST https://userstyles.world/oauth/access_token
```

#### Parameters

| Name |  Type | Required | Description |
| :--- | :--- | :--- | :--- |
| client_id | string | Yes | The client ID that is generated when you've registered for an USw application. |
| client_secret | string | Yes | The client secret that is generated when you've registered for an USw application. |
| code | string | Yes | The code that you've received as a parameter to Step 1.
| state | string | No | The unguessable random string you've provided in Step 1. |


#### Response

Depending on the `Accept` header that was sent.  
It will by default, send the `text/plain` format.

```
access_token=usw_IamAveryRandomAccesTokkeennAndImNotAfraid&token_type=bearer
```

Otherwise when specified:
```json
Accept: application/json
{"access_token":"usw_IamAveryRandomAccesTokkeennAndImNotAfraid", "token_type":"bearer"}
```

### 3 >> Using the access token

The access token will allow you to make requests to endpoints on behalf of the user.
The endpoints are specified [[./endpoints.MD][here]].

```
Authorization: TOKEN_TYPE ACCESS_TOKEN
GET https://userstyles.world/api/user/
```

## Scopes

_**Note:** Scopes may, later on, be split down into more specific scopes._

| Name | Description |
| :--- | :--- |
| `style` | Allow to add/edit/delete styles of the user. |
| `user` | Allow retrieving information of the user. |
