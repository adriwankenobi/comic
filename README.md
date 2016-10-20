# MARVEL comics management application

Get MARVEL API key pair from:

	http://developer.marvel.com/

## Usages

### (1) Create folders structure on hard drive

	go run main.go -folders -f marvel.xlsx -o <target-folder>
	
### (2) Update xlsx file with data from MARVEL API

	go run main.go -update -f marvel.xlsx -mpubkey <marvel_pub_key> -mprikey <marvel_private_key> -start 1998 -end 2016

### (3) Convert xslx file into different json files

	go run main.go -convert -f marvel.xlsx -o web/

### (4a) Deploy to local server

    cd web; goapp serve 

### (4b) Deploy to GAE

    cd web; appcfg.py -A <GAE_project_id> -V <version> update .
    
### (5) Test application

	localhost:8080
	http://<project_id>.appspot.com
	
### (6) Test queries

	curl -XGET -i localhost:8080/api/comics/:id
	curl -XGET -i http://<project_id>.appspot.com/api/comics/:id