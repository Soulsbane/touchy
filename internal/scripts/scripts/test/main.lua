print("Hello from main.lua")
print("Output Directory:", GetOutputDir())
Templates:ShowTemplate("go", "default")
Templates:List("all")
print("Has language Go: ", Templates:HasLanguage("go"))
print("Has language template go/default: ", Templates:HasTemplate("go", "default"))
print("Has language gozer: ", Templates:HasLanguage("gozer"))
print("Has language template gozer/busted: ", Templates:HasTemplate("gozer", "busted"))

