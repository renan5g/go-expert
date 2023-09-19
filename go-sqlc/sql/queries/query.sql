-- name: ListCategories :many
SELECT * FROM categories;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ?;

-- name: CreateCategory :exec
INSERT INTO categories (id, name, description)
VALUES (?, ?, ?);

-- name: UpdateCategory :exec
UPDATE categories SET name = ?, description = ?
WHERE id = ?;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = ?;

-- name: FindCategoryByCourseId :one
SELECT c.id, c.name, c.description 
FROM categories c JOIN courses co ON c.id = co.category_id 
WHERE co.id = ?;

-- name: CreateCourse :exec
INSERT INTO courses (id, name, description, price, category_id) 
VALUES (?, ?, ?, ?, ?);

-- name: ListCourses :many
SELECT c.*, ca.name AS category_name
FROM courses c JOIN categories ca ON c.category_id = ca.id;