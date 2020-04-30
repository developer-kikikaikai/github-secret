## Overview
This is a repository to set secret key by using GitHub API.  
Please get token from [Personal access tokens](https://github.com/settings/tokens) to use GitHub API

## Usage
```
./github-secret <command type> [options]
	<command type>, command type is one of following commands:update delete get list 
```
### ./github-secret update

Please call `make install` at first to install libsodium

```
  -owner string
        [required] repository owner (default "Please set the value")
  -repo string
        [required] repository name to access (default "Please set the value")
  -secname string
        [required] secret name (default "Please set the value")
  -secret string
        [required] secret value (default "Please set the value")
  -token string
    	access token (default environment value of "GITHUB_TOKEN")
```

### ./github-secret delete
```
  -owner string
        [required] repository owner (default "Please set the value")
  -repo string
        [required] repository name to access (default "Please set the value")
  -secname string
        [required] secret name (default "Please set the value")
  -token string
    	access token (default environment value of "GITHUB_TOKEN")
```
### ./github-secret get
```
  -owner string
        [required] repository owner (default "Please set the value")
  -repo string
        [required] repository name to access (default "Please set the value")
  -secname string
        [required] secret name (default "Please set the value")
  -token string
    	access token (default environment value of "GITHUB_TOKEN")
```
### ./github-secret list
```
  -owner string
        [required] repository owner (default "Please set the value")
  -repo string
        [required] repository name to access (default "Please set the value")
  -token string
    	access token (default environment value of "GITHUB_TOKEN")
```
