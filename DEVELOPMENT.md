# sauron

## Instalation

```
$ go get -u github.com/mguzelevich/sauron/...
```

## build

requirements
- nodejs
- `go get github.com/GeertJohan/go.rice/rice`
- ...

```
$ cd $ROOT/ui/frontend
$ npm install
$ npm run build

$ cd $ROOT/ui/
$ rice embed-go
```

## ui development

```
$ cd $ROOT/ui/frontend
$ npm install

# serve with hot reload at localhost:8080.
$ npm run dev

# build for production with minification
$ npm run build

# build for production and view the bundle analyzer report.
$ npm run build --report

# generate binary assets
$ cd $ROOT/ui/ && rice embed-go
```