package asgardian

func CreateSecurityMiddleware(strategy SecurityStrategy) SecurityMiddleware {
  return SecurityMiddleware{strategy}
}

func (app *Application) ConfigureSecurityMiddleware() SecurityMiddleware {
  strategy := CreateApiKeySecurityStrategy(app.config.Key)
  return CreateSecurityMiddleware(strategy)
}

