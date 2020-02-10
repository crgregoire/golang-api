# Buddha

## The source of _enlightenment_

## What Buddha is

Buddha is our core module for use with anything that would need to access our API. This includes AWS IOT Core, our eCommerce platform, our Interop dashboard, and more. It is the mediating precense between data in and out.

## Developer Setup

### Before you start

You will need [Git](https://git-scm.com/), [MySQL 5.7 (or later)](https://dev.mysql.com/downloads/), and [Golang 1.12 or later](https://golang.org/dl/). It is also recommended that you use [GoLand](https://www.jetbrains.com/go/) or [Visual Studio Code](https://code.visualstudio.com) for development.

Once you have those tools installed you should clone the repo [here](https://github.com/tespo/buddha.git).

### Setting up Buddha

To begin, open up the IDE of your choice (preferably Visual Studio Code, as GoLand requires a module import in the settings when you clone from GitHub) and use the terminal to pull the code using:

``` git
git pull
```

This ensures you have the latest code available.

After pulling run `make init`.  This will setup a couple different tools we've configured for this repository.

To begin, get the files `docker.env`, `local.env`, and `docker-compose.yml` from one of your teammates. When you've received those, place the `docker.env & local.env` files into a folder named `config` in the root directory of your project.

For the database, which will work closely with the tables that ***Satya*** creates, you must have an instance of MySQL running. Name the database `tespo_docker` (per the `docker.env`). You can do this by `docker-compose up -d mysqldb`. This requires the correct `docker-compose.yml` file in your project's root directory which you'll get from a teammate.

If you followed the directions to the letter, you should now direct your terminal instance to your working directory if you haven't already using:

``` bash
cd path/to/your/buddha's/root
```

When you're in this directory you can use the command:

``` golang
go build
```

This builds Buddha so you can test routes and their links to Satya.
The command to run Buddha after build is:

``` bash
./buddha
```

All tests must pass before you can merge `your-branch` into `develop`, then `develop` into `staging`, then `staging` into `master`.
