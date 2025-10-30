## Mongo File Upload Endpoints

Base path: `/go-fiber-mongo`

- Auth: Bearer JWT required. Admin can upload for any user; user can only for their own `:id`.

### Upload Photo
- Method: POST
- URL: `/go-fiber-mongo/users/:id/upload/photo`
- Body: `multipart/form-data` with field `file`
- Constraints: jpeg/jpg/png, max 1MB

### Upload Certificate
- Method: POST
- URL: `/go-fiber-mongo/users/:id/upload/certificate`
- Body: `multipart/form-data` with field `file`
- Constraints: pdf, max 2MB

Uploaded files are served statically at `/uploads/*`.


