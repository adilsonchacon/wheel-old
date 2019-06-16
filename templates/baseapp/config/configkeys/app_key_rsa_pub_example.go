package configkeys

var RsaPubExamplePath = []string{"config", "keys", "app.key.rsa.pub.example"}

var RsaPubExampleContent = `openssl rsa -in app.key.rsa -pubout > app.key.rsa.pub`
