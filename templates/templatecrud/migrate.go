package templatecrud

var MigrateContent = `model.Db.AutoMigrate(&entities.{{ .EntityName.CamelCase }}{})`
