print("Hello from main.lua")
print("Output Directory:", GetOutputDir())
Templates:ShowTemplate("go", "default")
Templates:List("all")
print("Has language Go: ", Templates:HasLanguage("go"))
print("Has language brainfuck: ", Templates:HasLanguage("brainfuck"))

