# Agile-Software-Development
## Assignment Description:
Create an application that can help patients around the US to find the cheapest option for treatment. This application will be able to show the user these options in a list format and on a map. The user should be able to sort by price, location and a "best match".

## Live version
http://99.81.88.54/index.html

## Links to videos on YouTube 
* Scrum meeting - https://youtu.be/zfyHyfF23tk
* Pair programming one - https://youtu.be/ZJvZb5cO-xs
* Pair programming two - https://youtu.be/w8D8sQ1oQzE
* TDD video evidnce - https://www.youtube.com/watch?v=fN_Gcm8Ta2M&feature=youtu.be


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

Required extensions:
```
CREATE EXTENSION cube;
CREATE EXTENSION earthdistance;
```

In order to migrate to latest version of DB schema, run:
```
go run main.go -migrate
```  
