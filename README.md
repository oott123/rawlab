# RawLab

Serves raw files from gitlab, with correct MIME types.

## Run

```bash
export GITLAB_API=https://gitlab.com/api/v4
export PORT=8888
./rawlab
```

or with docker:

```bash
docker run --rm -p 8625:8625 -e GITLAB_API=https://gitlab.com/api/v4 quay.io/oott123/rawlab:master
```

## Get

### Without Authorization

```http request
GET /username/repo@branch/path/to/file.ext HTTP/1.1
```
### With Token in Header

```http request
GET /username/repo@branch/path/to/file.ext HTTP/1.1

Authorization: Bearer glpat-***************
```

### With Token in URL

```http request
GET /username/repo@branch/path/to/file.ext?token=glpat-*************** HTTP/1.1
```

## Why

* to serve deno packages
* to serve simple static web sites

## Reminds

* this will serve HTML files as text/html, be sure to trust your users!