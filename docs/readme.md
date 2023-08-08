<b>build</b> 
- go run github.com/99designs/gqlgen generate --verbose
- docker compose --env-file ./configs/development.env up -d
<b>run</b>
- go run main.go development | production mode

<b>questions & issues</b>
- При удалении проекта, необходимо удалять ассеты. на ассет сервере нужно создать таблицу имени файла и id проекта.
  И создать ручку удаления ассетов по id проекта
- Может ли unit admin активировать пользователей?
- Кто кого из пользователей может обновлять, создавать, удалять?
- Сразу ли активен пользователь после создания через CreateUser?
- При удалении пользователя надо удалять все связи (projectPage, project)