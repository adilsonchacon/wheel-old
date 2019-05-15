# Wheel

Wheel is a framework written in Go Language for developers to quickly and easily build APIs and microservices.

Wheel follows the Model-View-Controller (MVC) architectural pattern. By default, it uses JWT for authentication, Postgresql and responses JSON data. But you can modify any of them, if you want to.

## Table of contents

-  [Install](#install)
-  [Set Up](#set-up)
	- [JWT](#jwt)
	- [Database](#database)
	- [Email](#email)
	- [Application](#application)
-  [Running](#set-up)
-  [File System](#file-system)
-  [Features](#features)
	- [Routes](#routes)
	- [Middleware](#middleware)
	- [Authentication](#authentication)
	- [Authorization](#authorization)
	- [ORM](#orm)
	- [Migration](#migration)
	- [Search Engine](#search-engine)
	- [Pagination](#pagination)
	- [Internacionalization](#internacionalization)
	- [Sending Emails](#sending-emails)
	- [Logging](#logging)
-  [Programming](#programming)
	- [Adding Single Service](#adding-single-service)
	- [Adding New CRUD](#adding-new-crud)
-  [License](#license)

## Install

### Go

[Install Golang](https://golang.org/doc/install)

### Dependences

```
$> go get github.com/jinzhu/gorm
$> go get gopkg.in/yaml.v2
$> go get github.com/gorilla/mux
$> go get github.com/dgrijalva/jwt-go
$> go get github.com/satori/go.uuid
$> go get github.com/lib/pq
$> go get golang.org/x/crypto/bcrypt
```


## Set up


### Step 1: JWT

The following keys are used by JWT.

```
$> cd config/keys
$> openssl genrsa -out app.key.rsa 1024
$> openssl rsa -in app.key.rsa -pubout > app.key.rsa.pub
```

### Step 2: Database

Connect to your database

```
$> cd config
$> cp database.example.yml database.yml
```

Edit _database.yml_ and set up your database credentials

### Step 3: Email

Connect to your email provider

```
$> cd config
$> cp email.example.yml email.yml
```

Edit _email.yml_ and set up your email credentials

### Step 4: Application

Set up _app.yml_ file for app work pristine

```
$> cd config
$> cp app.example.yml app.yml
```

Edit app.yml and set the following options:


| Item | Definition |
| ------ | ----------- |
| _app_name_ | Give a name to your app |
| _frontend_base_url_ | URL to be used on your frontend |
| _secret_key_ | Key to encrypt passwords on database |
| _reset_password_expiration_seconds_ | After reset password, how long (in seconds) is it valid? |
| _token_expiration_seconds_ | After a JWT token is generated, how long (in seconds) is it valid? |
| _locales_ | List of available locales |

### Step 5: Locales

Words and phrases for internacionalization. You can add your own locales files, but remenber to set up the _app.yml_ configuration file.

```
$> cd config/locales
$> cp en.yml YOUR_LOCALE.yml
```

And edit _YOUR_LOCALE.yml_


## Running

Before running you should be sure your database schema is up to date, just run _migrate_ mode:

```
$> go run main.go -mode=migrate
```

Running:

```
$> go run main.go
```

### Help

Prints help menu

```
$> go run main.go --help
```

Will print:

```
  -mode string
        run mode (options: server/migrate) (default "server")
  -host string
        http server host (default "localhost")
  -port string
        http server port (default "8081")
```


You can set up the Port and Host binding



## File System

![File System](http://static.smart26.com/file_system.png)]

### App/Controllers
---

By Default, _Wheel_ has three controllers: Users, Sessions and Myself.

| Controller Name | Definition |
| ------ | ----------- |
| Users | Handles entire CRUD for users. Only _admin_ users have access to this CRUD. |
| Sessions | Handles sing up, sign in, sign out and password recovery functionalities. |
| Myself | Users can modify their own data and password. |


#### Note

> The file with _app_ prefix is not a controller, but it has methods those will be used by all controllers.

### App/Models
---
By Default, _Wheel_ has two models: User and Session.

| Model Name | Definition |
| ------ | ----------- |
| User | It persists and extracts data in the database. Table: _users_ |
| Session | It persists and extracts data in the database. Table: _sessions_ |

Relationship:

```
One user has many sessions
```


#### Note

> Files with _app_ prefix are not models, but it has methods those will be used by all models. Each file with prefix _app_ has specifics functionalities:

| File Name | Definition |
| ------ | ----------- |
| app_migration.go | User for database schema migration. |
| app_model.go | Handles database conection and others helpfull methods. |
| app_paginate.go | Handles pagination for lists. |
| app_search_engine.go | Handles native search for list. |


### App/Views
---

By Default, _Wheel_ has three view: User, Session and Myself.

| View Name | Definition |
| ------ | ----------- |
| User | Notifies, lists and shows users data. |
| Session | Session's notifications. |
| Myself | Notifies and shows personal data. |

#### Note

> The file with _app_ prefix is not a view, but it has methods those will be used by all views.


#### App/Views/Mailer

E-mails templates for Password Recovery and Sign Up Confirmation.

### Config
---

Here you can find the configuration files.

### Log
---

Log repository.

### Routes
---

The file _routes.go_ sets up the API interface and the file _middleware.go_ is the pre-processor of the requests.

### Utils
---

Utilities to be used in all framework.

| File | Definition |
| ------ | ----------- |
| converter | Functions for type convertion |
| locale | Functions for locale |
| logger | Functions for logging |
| mailer | Functions for sending mail |
| safer | Functions for cryptography |


## Features


### Routes

Routes are RESTful. Each entry has one HTTP verb, one URL path and one controller's method whitch responds the requests.

### Middleware

_Middlewares_ are executed before a request reachs the controller's method. By default _Middleware_ handles the _Authentication_ and _Authorization_. See bellow both.

### Authentication



### Authorization

### ORM

### Migration

### Search Engine

### Pagination

### Internacionalization

### Sending Emails

### Logging

log.Info.Println("Logger example")

## Programming


### Adding Single Service

### Adding New CRUD


## License
