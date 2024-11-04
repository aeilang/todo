-- name: GetTodos :many
select id, content, completed, create_at from todos where is_deleted = false order by id;

-- name: UpdateTodoCompleted :exec
update todos set completed = not completed where id = $1 and is_deleted = false;

-- name: DeleteTodo :exec
update todos set is_deleted = true where id = $1;

-- name: CreateTodo :exec
insert into todos(content) values ($1);