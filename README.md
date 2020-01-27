# Agile-Software-Development
## Assignment Description:
Create an application that can help patients around the US to find the cheapest option for treatment. This application will be able to show the user these options in a list format and on a map. The user should be able to sort by price, location and a best match,

## Project Information:
The project utilises the following technologies.
* Frontend:
	* HTML
	* CSS
	* JavaSctipt / JQuery
* Backend:
	* Go 1.13
* Database:
	* PostgreSQL 11.5+

## TODO:
Currently tracked via trello.

## DONE
Currently tracked via trello.


## PostgreSQL setup for local dev env:
For dev environment user `cms` with password `secret` is required  
For unit testing user `test` with password `test` is required  
Once the users are setup, create databases `cms` and `test_cms`  
In order to migrate to latest version of DB schema, run:
```
go run main.go -migrate
```  
