package main

type Formatter interface {
	Format() (formatted string)
}

type PlainText struct {
	message string
}

// Notice that we don't need to name the argument here again
func (p PlainText) Format() string {
	return p.message
}

type Bold struct {
	message string
}

// But we can still name it if we want
func (b Bold) Format() string {
	return "**" + b.message + "**"
}

type Code struct {
	message string
}

// Doesn't even have to be the same name
func (c Code) Format() (inlineCode string) {
	return "`" + c.message + "`"
}

// Don't Touch below this line

func SendMessage(formatter Formatter) string {
	return formatter.Format() // Adjusted to call Format without an argument
}
