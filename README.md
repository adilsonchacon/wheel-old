## Overview

The Wheel framework is a single and executable file. It runs through command line and generates codes (in Go Language) for you to build your RESTful API, avoiding rework while building the API architecture and it maintenance. 

## Features

- [MVC](http://wheel.unity26.org/features#mvc)
- [RESTful](http://wheel.unity26.org/features#restful)
- [JWT](http://wheel.unity26.org/features#jwt)
- [Session controller](http://wheel.unity26.org/features#session-controller)
- [Middleware](http://wheel.unity26.org/features#middleware)
- [Authorization](http://wheel.unity26.org/features#authorization)
- [Users management](http://wheel.unity26.org/features#users-management)
- [ORM](http://wheel.unity26.org/features#orm)
- [Search engine](http://wheel.unity26.org/features#search-engine)
- [Lists pagination](http://wheel.unity26.org/features#lists-pagination)
- [I18n](http://wheel.unity26.org/features#i18n)
- [Sends email](http://wheel.unity26.org/features#sends-email)
- [Log](http://wheel.unity26.org/features#log)


See full documentation of default features at http://wheel.unity26.org/features


## Install

### Go

[Install Golang](https://golang.org/doc/install)

### Dependences

```
$> go get github.com/iancoleman/strcase
$> go get github.com/jinzhu/inflection
```

### Wheel

```
$> go get github.com/adilsonchacon/wheel
$> cd GOPATH/src/github.com/adilsonchacon/wheel
$> go build -o wheel main.go 
$> sudo mv wheel /usr/bin
```

__GOPATH__ is where the Go packages and sources are installed

At this example the executable file was moved to _/usr/bin_. But feel free to set it up to any directory you want. Just add the path to your _.profile_, as you see below:

```
export PATH=$PATH:YOUR\_DESIRED\_PATH
```


## Usage

Wheel has basically two options: _new_ to create new APIs and _generate_ to add new functionalities to your API. 

Check _help_ for more details.


```
wheel --help
```

### New Application 

Let's create an API for a Blog.

```
wheel new github.com/account_name/blog
```

It will output something like this:

```
"Go" seems installed
Checking dependences...
         ...
Generating new app...
         created: GOPATH/src/github.com/account_name/blog
         ...

Your App was successfully created!

Change to the root directory using the command line below: 
cd GOPATH/src/github.com/account_name/blog

Set up your database connection modifying the file config/database.yml

For more details call help:
go run main.go --help
```

Remember: __GOPATH__ is where the Go packages and sources are installed


### Configure 

####  Database

Currently, Wheel has support only for Postgresql. Edit _config/database.yml_ and set up your database connection.

#### Email

To connect to your email provider edit _config/email.yml_ and set up with your send email account.

#### Application

Edit _config/app.yml_ and set the following options:


| Item | Definition |
| ------ | ----------- |
| _app_name_ | Your app name |
| _app_repository_ | Repository name |
| _frontend_base_url_ | URL to be used on your frontend |
| _secret_key_ | Key to encrypt passwords on database |
| _reset_password_expiration_seconds_ | After reset password, how long (in seconds) is it valid? |
| _token_expiration_seconds_ | After a JWT token is generated, how long (in seconds) is it valid? |
| _locales_ | List of available locales |

#### Locales

Words and phrases for internacionalization. You can add your own locales files, but remember to add to _config/app.yml_ configuration file first.


### Running

Before running you should be sure your database schema is up to date, just run the _migrate_ mode:

```
$> go run main.go -mode=migrate
```

Run:

```
$> go run main.go
```

#### Checking

Before continuing, please, check file _db/schema/migrate.go_ and get user's email and password.

Below is part of full documentation of default resources. See full documentation at http://wheel.unity26.org/default-resources

__POST /sessions/sign\_in__

| Form param | Description | Mandatory |
| ------ | ----------- | ----------- |
| email | User's email | Yes |
| password | User's password | Yes |


On success:

```
{
  "system_message": {
      "type": "notice",
      "content": "signed in successfully"
  },
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9....",
  "expires": 7200
}
```

On error:

```
{
  "system_message": {
      "type": "alert",
      "content": "could not sign in"
  },
  "errors": ["invalid credentials"]
}
```

__DELETE /sessions/sign\_out__

| Header param | Description | Mandatory |
| ------ | ----------- | ----------- |
| token | Session token | Yes |


On success:

```
{
  "system_message": {
    "type":"notice",
    "content":"signed out successfully"
  }
}
```

On error:

```
{
  "system_message": {
      "type": "alert",
      "content": "access denied"
  },
  "errors": ["invalid token"]
}
```

See full documentation of default resources at http://wheel.unity26.org/default-resources


### New CRUD

Based on the Blog API above, let's create a new CRUD.

Don't forget to call the directory where the application were generated.
 
```
cd GOPATH/src/github.com/account_name/blog
wheel g scaffold post title:string description:text published:bool user:references
```

It will output something like this:

```
"Go" seems installed
Checking dependences...
         _..._
Generating new CRUD...
         created: app/post/post_model.go
         created: app/post/post_view.go
         created: db/entities/post_entity.go
         created: app/handlers/post_handler.go
         updated: routes/routes.go
         updated: db/schema/migrate.go
         updated: routes/authorize.go
```

After any changing, don't forget to run the _migrate_ mode:

```
$> go run main.go -mode=migrate
$> go run main.go
```

## Full Documentation

See full documentation at http://wheel.unity26.org/


## License

Wheel is released under the [MIT License](https://opensource.org/licenses/MIT).