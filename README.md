# Snippet Box

> [!NOTE]
> How to know my docker container gateway for connection?  
> `docker inspect -f '{{range .NetworkSettings.Networks}}{{.Gateway}}{{end}}' <container_name_or_id>`

To login to mysql database:

`mysql -h <your_database_gateway> -u web -p <your_database_name>`
