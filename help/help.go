package help

import ()

var Content = `
Usage:
  wheel new APP_PATH [options]             # Creates new app

  wheel generate SUBJECT NAME COLUMNS      # Adds new CRUD to an existing app. 
                                           # SUBJECT options: scaffold/model/entity/view/handler. 
                                           # NAME: scaffold, mode, entity, view or handler name
                                           # COLUMNS: pair of column name and column type separated by ":"
                                           # i.e. description:string
                                           # Available types are: string/text/integer/decimal/datetime/bool/references
                                           
Options:
  -G, [--skip-git]                         # Skip .gitignore file

More:
  -h, [--help]                             # Show this help message and quit
  -v, [--version]                          # Show Wheel version number and quit
`
