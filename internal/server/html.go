package server

import (
	"fmt"
	"strings"

	"github.com/zachmann/offa/internal/config"
)

const _loginHtml = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="$BASEPATH$/static/css/sakura-earthly.css" media="screen" />
<link rel="stylesheet" href="$BASEPATH$/static/css/sakura-dark.css" media="screen and (prefers-color-scheme: dark)" />
<link rel="stylesheet" href="$BASEPATH$/static/css/offa.css" />
</head>
<body>
<h2>%s</h2>
%s
<h4>Choose an OP to login</h4>
<form action="">
  <select name="op" id="op">
	<option value="/" selected disabled>Choose OP...</option>
	%s
  </select>
  <input type="hidden" name="next" value="%s" />
  <input type="submit" value="Login">
</form>
</body>
</html>`

const _errorHtml = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="$BASEPATH$/static/css/sakura-earthly.css" media="screen" />
<link rel="stylesheet" href="$BASEPATH$/static/css/sakura-dark.css" media="screen and (prefers-color-scheme: dark)" />
<link rel="stylesheet" href="$BASEPATH$/static/css/offa.css" />
</head>
<body>
<h3>Error %s</h3>
<p>%s</p>
<a href="%s">Back to Login</a>
</body>
</html>`

var errorHtml string
var loginHtml string

func initHtmls() {
	errorHtml = strings.ReplaceAll(_errorHtml, "$BASEPATH$", config.Get().Server.Basepath)
	loginHtml = strings.ReplaceAll(_loginHtml, "$BASEPATH$", config.Get().Server.Basepath)
}

func errorPage(error, message string) string {
	return fmt.Sprintf(errorHtml, error, message, getFullPath(config.Get().Server.Paths.Login))
}
